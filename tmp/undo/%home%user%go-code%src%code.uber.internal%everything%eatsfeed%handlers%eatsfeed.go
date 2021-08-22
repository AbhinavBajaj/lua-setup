Vim�UnDo� �$dۉ�;W2AUKT�px�	E|r��K�Q.�  X                                   a�3     _�                      �   L    ����                                                                                                                                                                                                                                                                                                                                                             a�2     �              X   package handlers       import (   
	"context"   		"errors"   	"math/rand"   
	"strconv"   	"time"       E	eatsfeedt "thriftrw/code.uber.internal/everything/eatsfeed/eatsfeed"       &	"code.uber.internal/eats/base/safego"   2	"code.uber.internal/everything/eatsfeed/business"   4	"code.uber.internal/everything/eatsfeed/dependency"   -	"code.uber.internal/everything/eatsfeed/lib"   ;	"code.uber.internal/everything/eatsfeed/lib/jaegerhelpers"   9	"code.uber.internal/everything/instrumenter.git/factory"   /	"code.uber.internal/everything/loadshedder/v2"   )	"code.uber.internal/go-common.git/x/log"   2	"code.uber.internal/rt/flipr-client-go.git/flipr"   :	"code.uber.internal/rt/go-common.git/handler/interceptor"   !	"github.com/andres-erbsen/clock"   6	apacheThrift "github.com/apache/thrift/lib/go/thrift"   	"github.com/uber/tchannel-go"   %	"github.com/uber/tchannel-go/thrift"   )       const (   ;	enableFallbackCachedFeedFlipr = "eater.feed.fallbackCache"   A	forceFallbackCachedFeedFlipr  = "eater.feed.fallbackCache.force"   5	getFeedTimeoutFlipr           = "eater.feed.timeout"       @	fallbackCacheUserTriggeredKey  = "fallbackCache.user.triggered"   2	getFeedStatsM3Key              = "get_feed.stats"   >	getConstrainedFeedStatsM3Key   = "get_constrained_feed.stats"   6	getFeedTimeoutKey              = "eater.feed.timeout"   @	enableFallbackWhenTimeoutFlipr = "fallback.when.timeout.enable"       u	// GetFeed requests in RTAPI times out in 7s, the default value is close enough to the RTAPI timeout and also leaves   /	// some buffer for sending data over the wires   n	// TODO(zhaokun): revisit the value after the feed re-architecture project is rolled out, we might be able to   	// decrease it   1	getFeedTimeoutFallback = 6800 * time.Millisecond       (	fallbackFeedLoggingTag = "fallbackFeed"    	feedLoggingTag         = "feed"   )       0// FeedHandler is used encapsulate feed handlers   type FeedHandler struct {   	Deps        *dependency.Deps   *	FeedService business.FeedServiceInterface   	Rand        *rand.Rand   }       /// NewFeedHandler is an instance of FeedHandler   Gfunc NewFeedHandler(deps *dependency.Deps) eatsfeedt.TChanFeedService {   	serviceImpl := &FeedHandler{   		Deps:        deps,   4		FeedService: business.InitializeFeedService(deps),   ?		Rand:        rand.New(rand.NewSource(time.Now().UnixNano())),   	}       Q	// Only loadshed feed service; recommendation service, while called via muttley,   G	// is actually only called from eatsfeed itself, so this is "upstream"   C	ls := loadshedder.NewCodelInterceptorOrPanic(context.Background(),   !		&loadshedder.InterceptorConfig{   			ServiceName: "eatsfeed",   			ErrorFunc: func() error {   ?				return &eatsfeedt.InternalError{Info: &eatsfeedt.ErrorInfo{   .					ShouldRetry: apacheThrift.BoolPtr(false),   ?					Message:     apacheThrift.StringPtr("LoadShedding error"),   				}}   			},   		},   4		&loadshedder.LoadSheddingConfig{Name: "eatsfeed"},   		deps.Logger(),   		deps.Metrics(),   		deps.Flipr(),   	)       i	interceptor := eatsfeedt.FeedServiceInterceptor{GlobalInterceptors: []interceptor.GlobalInterceptor{ls}}   4	toThriftErr := func(err error) error { return err }       e	return eatsfeedt.NewFeedServiceRequestHandler(serviceImpl, interceptor, deps.Metrics(), toThriftErr)   }       // GetFeed returns a feed   t// workflow: when the fallback cache is enabled, we always build the fallback feed in a concurrent go-routine at the   m// beginning of the request. if the regular feed or the fast food cached feed doesn't return a non-empty feed   :// within the timeout range, we'd return the fallback feed   ofunc (fh FeedHandler) GetFeed(ctx thrift.Context, request *eatsfeedt.GetFeedRequest) (*eatsfeedt.Feed, error) {   >	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetFeed")   	if span != nil {   		defer span.Finish()   	}       \	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertFeedParams(newCtx, request)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   "		}).Error("Invalid feed request")   		return nil, err   	}       	var feed eatsfeedt.Feed    	var fallbackFeed eatsfeedt.Feed   p	defaultConstraints := flipr.Constraints{flipr.Constraints2(lib.GetDefaultFliprConstraints2Map(*getFeedParams))}   6	// whether to allow users to access the fallback feed   v	enableFallbackCache, _ := fh.Deps.Flipr().GetBoolValue(ctx, enableFallbackCachedFeedFlipr, defaultConstraints, false)   3	// whether to force users to use the fallback feed   s	forceFallbackFeed, _ := fh.Deps.Flipr().GetBoolValue(ctx, forceFallbackCachedFeedFlipr, defaultConstraints, false)   }	enableFallbackWhenTimeout, _ := fh.Deps.Flipr().GetBoolValue(ctx, enableFallbackWhenTimeoutFlipr, defaultConstraints, false)       +	errChannelFallbackFeed := make(chan error)   	if enableFallbackCache {   O		// this will not degrade feed performance since it's happening asynchronously   q		// TODO(zhaokun): 1. after running and monitoring the fallback feed in production for a few months, remove this   r		// feature flag  2. fallback feed will be migrated to a higher layer(eater-gateway possibly) after the west plan   		// migration   <		errChannelFallbackFeed = fh.Deps.Instrumenter.AsyncInvoke(   			ctx,   			fh.Deps.ZapLogger(),   			"loadFallbackFeed",   $			func(ld *factory.LogData) error {       				var fallbackFeedErr error   Z				fallbackFeed, fallbackFeedErr = fh.FeedService.GetFallbackFeed(newCtx, *getFeedParams)   				return fallbackFeedErr   			},   		)       		if forceFallbackFeed {   !			err = <-errChannelFallbackFeed   ?			fh.logFeedErrors(err, request, fallbackFeedLoggingTag, true)   			return &fallbackFeed, err   		}   		} else {   		close(errChannelFallbackFeed)   	}       E	getFeedFunc := func(reqCtx thrift.Context) (eatsfeedt.Feed, error) {   =		feed, err := fh.FeedService.GetFeed(reqCtx, *getFeedParams)   		return feed, err   	}       6	if enableFallbackCache && enableFallbackWhenTimeout {   M		feed, err = fh.getFeedWithTimeout(getFeedFunc, ctx, fh.getFeedTimeout(ctx))   		} else {   !		feed, err = getFeedFunc(newCtx)   	}       5	if enableFallbackCache && len(feed.FeedItems) == 0 {    		err = <-errChannelFallbackFeed   		feed = fallbackFeed   >		fh.logFeedErrors(err, request, fallbackFeedLoggingTag, true)   		} else {   7		fh.logFeedErrors(err, request, feedLoggingTag, false)   	}   	return &feed, err   }       "// GetMapFeed generates a map feed   func (fh FeedHandler) GetMapFeed(ctx thrift.Context, request *eatsfeedt.GetFeedRequest) (response *eatsfeedt.Feed, err error) {   "	defer fh.Deps.Instrumenter.Start(   		ctx,   		fh.Deps.ZapLogger(),   		"Eatsfeed.GetMapFeed",   :	).Finish(func(ctx context.Context, ld *factory.LogData) {   		lib.LogErrorMsg(err, ld)       	})       >	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetFeed")   	if span != nil {   		defer span.Finish()   	}       \	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertFeedParams(newCtx, request)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   "		}).Error("Invalid feed request")   		return nil, err   	}       <	feed, err := fh.FeedService.GetMapFeed(ctx, *getFeedParams)   	if err != nil {   		return nil, err   	}       	return &feed, nil   }       $func (fh FeedHandler) logFeedErrors(   	err error,   #	request *eatsfeedt.GetFeedRequest,   	feedPath string,   	isFallbackFeed bool,   ) {   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":    err.Error(),   			"request":  request,   			"feedPath": feedPath,   		}).Error("Eatsfeed.GetFeed")   	}       	if isFallbackFeed {   H		fh.Deps.Metrics().Counter(fallbackCacheUserTriggeredKey).Inc(int64(1))   	}       5	fh.Deps.Metrics().Counter(getFeedStatsM3Key).Tagged(   		map[string]string{   5			"dining_mode":   request.GetDiningMode().String(),   7			"fallback_feed": strconv.FormatBool(isFallbackFeed),   		},   	).Inc(int64(1))   }       '// set a timeout on the GetFeed request   �func (fh FeedHandler) getFeedWithTimeout(getsFeed func(ctx thrift.Context) (eatsfeedt.Feed, error), ctx thrift.Context, timeout time.Duration) (eatsfeedt.Feed, error) {   	type GetFeedResult struct {   		feed eatsfeedt.Feed   		err  error   	}   #	cnl := make(chan GetFeedResult, 1)   	safego.Go(func() {   Y		getFeedWithTimeoutCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetFeedWithTimeout")   		if span != nil {   			defer span.Finish()   		}   O		// set the same timeout on the request context - cancel requests upon timeout   K		timeoutCtx, cancel := context.WithTimeout(getFeedWithTimeoutCtx, timeout)   p		tTimeoutCtx := tchannel.WrapWithHeaders(timeoutCtx, lib.CloneStringStringMap(getFeedWithTimeoutCtx.Headers()))   		defer cancel()       $		feed, err := getsFeed(tTimeoutCtx)   		result := new(GetFeedResult)   		result.feed = feed   		result.err = err   		cnl <- *result   	})    	timer := time.NewTimer(timeout)   	defer timer.Stop()   		select {   	case result := <-cnl:    		return result.feed, result.err   	case <-timer.C:   <		fh.Deps.Metrics().Counter(getFeedTimeoutKey).Inc(int64(1))   i		return eatsfeedt.Feed{}, errors.New("GetFeed request timed out, the fallback feed should be triggered")   	}   }       :// GetPickupFeed returns a feed for the pickup dining mode   �func (fh FeedHandler) GetPickupFeed(ctx thrift.Context, request *eatsfeedt.GetPickupFeedRequest) (*eatsfeedt.GetPickupFeedResponse, error) {   D	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetPickupFeed")   	if span != nil {   		defer span.Finish()   	}       c	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertPickupFeedRequest(newCtx, request)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   )		}).Error("Invalid pickup feed request")   		return nil, err   	}       ?	feed, err := fh.FeedService.GetPickupFeed(ctx, *getFeedParams)   	if err != nil {   R		fh.Deps.Logger().WithField("error", err.Error()).Error("Eatsfeed.GetPickupFeed")   	}       	return &feed, nil   }       B// GetScheduledFeed returns a feed for the early bird landing page   �func (fh FeedHandler) GetScheduledFeed(ctx thrift.Context, request *eatsfeedt.GetScheduledFeedRequest) (*eatsfeedt.GetScheduledFeedResponse, error) {   >	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetFeed")   	if span != nil {   		defer span.Finish()   	}       f	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertScheduledFeedRequest(newCtx, request)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   ,		}).Error("Invalid scheduled feed request")   		return nil, err   	}        	// Construct storeindex request   B	feed, err := fh.FeedService.GetScheduledFeed(ctx, *getFeedParams)   	if err != nil {   U		fh.Deps.Logger().WithField("error", err.Error()).Error("Eatsfeed.GetScheduledFeed")   	}       	return &feed, nil   }       *// GetCateringFeed returns a catering feed   &func (fh FeedHandler) GetCateringFeed(   	ctx thrift.Context,   +	request *eatsfeedt.GetCateringFeedRequest,   ) (   $	*eatsfeedt.GetCateringFeedResponse,   	error,   ) {   >	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetFeed")   	if span != nil {   		defer span.Finish()   	}       e	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertCateringFeedRequest(newCtx, request)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   +		}).Error("Invalid catering feed request")   		return nil, err   	}       A	feed, err := fh.FeedService.GetCateringFeed(ctx, *getFeedParams)   	if err != nil {   T		fh.Deps.Logger().WithField("error", err.Error()).Error("Eatsfeed.GetCateringFeed")   	}   	return &feed, err   }       A// GetConstrainedFeed returns feed after applying the constraints   lfunc (fh FeedHandler) GetConstrainedFeed(ctx thrift.Context, request *eatsfeedt.GetConstrainedFeedRequest) (   0	*eatsfeedt.GetConstrainedFeedResponse, error) {   I	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetConstrainedFeed")   	if span != nil {   		defer span.Finish()   	}       	fh.Deps.Metrics().Tagged(   		map[string]string{   G			"diningMode":     request.GetFeedRequest().GetDiningMode().String(),   X			"hasFilterInput": strconv.FormatBool(len(request.GetFeedRequest().GetFilters()) > 0),   		},   6	).Counter(getConstrainedFeedStatsM3Key).Inc(int64(1))   r	getConstrainedFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertGetConstrainedFeedParams(ctx, request)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   "		}).Error("Invalid feed request")   		return nil, err   	}       ]	constrainedFeed, err := fh.FeedService.GetConstrainedFeed(newCtx, *getConstrainedFeedParams)   	if err != nil {   W		fh.Deps.Logger().WithField("error", err.Error()).Error("Eatsfeed.GetConstrainedFeed")   	}   	return &constrainedFeed, err   }       -// GetValueFeed returns a value-oriented feed   �func (fh FeedHandler) GetValueFeed(ctx thrift.Context, request *eatsfeedt.GetValueFeedRequest) (*eatsfeedt.GetValueFeedResponse, error) {   >	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetFeed")   	if span != nil {   		defer span.Finish()   	}       a	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertValueFeedParams(newCtx, request)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   (		}).Error("Invalid Value Feed request")   		return nil, err   	}       A	feed, err := fh.FeedService.GetValueFeed(newCtx, *getFeedParams)   	if err != nil {   Q		fh.Deps.Logger().WithField("error", err.Error()).Error("Eatsfeed.GetValueFeed")   	}   	return &feed, err   }       0// GetFriendFeed returns a feed of friend orders   �func (fh FeedHandler) GetFriendFeed(ctx thrift.Context, request *eatsfeedt.GetFriendFeedRequest) (*eatsfeedt.GetFriendFeedResponse, error) {   D	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetFriendFeed")   	if span != nil {   		defer span.Finish()   	}       c	params, err := lib.NewRequestConstructor(fh.Deps).ConvertFeedParams(newCtx, request.FeedItemParam)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   )		}).Error("Invalid Friend Feed request")   		return nil, err   	}       T	response, err := fh.FeedService.GetFriendFeed(newCtx, *params, request.FriendUuids)   	if err != nil {   R		fh.Deps.Logger().WithField("error", err.Error()).Error("Eatsfeed.GetFriendFeed")   	}   	return response, err   }       F// GetCityStoresFeed returns a feed for getting all stores in the city   �func (fh FeedHandler) GetCityStoresFeed(ctx thrift.Context, request *eatsfeedt.GetFeedRequest, pageInfo *eatsfeedt.Pagination) (*eatsfeedt.Feed, error) {   Y	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertFeedParams(ctx, request)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   "		}).Error("Invalid feed request")   		return nil, err   	}       K	if pageInfo == nil || pageInfo.PageSize == nil || pageInfo.Offset == nil {   		fh.Deps.Logger().   *			WithField("err", "invalid Pagination").   			WithField("query", request).   *			Error("Pagination should be specified")   X		return nil, lib.WrapParamsError("InvalidPagination", errors.New("invalid Pagination"))   	}       N	feed, err := fh.FeedService.GetCityStoresFeed(ctx, *getFeedParams, *pageInfo)   	if err != nil {   V		fh.Deps.Logger().WithField("error", err.Error()).Error("Eatsfeed.GetCityStoresFeed")   	}   	return &feed, err   }       E// GetRestaurantRewardsFeed returns a feed for Restaurant Rewards Hub   �func (fh FeedHandler) GetRestaurantRewardsFeed(ctx thrift.Context, request *eatsfeedt.GetRestaurantRewardsFeedRequest) (*eatsfeedt.GetRestaurantRewardsFeedResponse, error) {   e	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertFeedParams(ctx, request.FeedRequest)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   6		}).Error("Invalid GetRestaurantRewardsFeed request")   		return nil, err   	}       D	return fh.FeedService.GetRestaurantRewardsFeed(ctx, *getFeedParams)   }       p// GetSimilarRecommendation returns a feed containing similar recommendation, e.g. similar stores of given store   �func (fh FeedHandler) GetSimilarRecommendation(ctx thrift.Context, request *eatsfeedt.GetSimilarRecommendationRequest) (*eatsfeedt.GetSimilarRecommendationResponse, error) {   G	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetSimilarStores")   	if span != nil {   		defer span.Finish()   	}       e	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertFeedParams(ctx, request.FeedRequest)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   .		}).Error("Invalid GetSimilarStores request")   		return nil, err   	}       a	// GetSimilarRecommendation can serve as a generic similar recommendation endpoint going forward   3	// for now, it only has the getSimilarStores logic   	if request.StoreUUID == nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"request": request,   ?		}).Error("Invalid GetSimilarStores request: Empty storeUUID")   		return nil, nil   	}       	pushBackToClient := false   B	if request.PushBackToClient != nil && *request.PushBackToClient {   		pushBackToClient = true   	}   �	resp, err := fh.FeedService.GetSimilarStores(newCtx, *getFeedParams, string(request.GetStoreUUID()), request.GetStoreName(), request.SeenStoreUUIDs, pushBackToClient)   	if err != nil {   U		fh.Deps.Logger().WithField("error", err.Error()).Error("Eatsfeed.GetSimilarStores")   	}   	return &resp, err   }       .// GetHealth returns the health of the service   Pfunc (fh FeedHandler) GetHealth(ctx thrift.Context) (*eatsfeedt.Health, error) {   	healthVal := true   .	return &eatsfeedt.Health{Ok: &healthVal}, nil   }       2// GetFeedUpdate returns the updates to feed items   �func (fh FeedHandler) GetFeedUpdate(ctx thrift.Context, request *eatsfeedt.GetFeedUpdateRequest) (*eatsfeedt.GetFeedUpdateResponse, error) {   D	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetFeedUpdate")   	if span != nil {   		defer span.Finish()   	}       c	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertFeedUpdateRequest(newCtx, request)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   +		}).Error("Invalid GetFeedUpdate request")   		return nil, err   	}       E	feedUpdate, err := fh.FeedService.GetFeedUpdate(ctx, *getFeedParams)   	if err != nil {   R		fh.Deps.Logger().WithField("error", err.Error()).Error("Eatsfeed.GetFeedUpdate")   	}       	return feedUpdate, nil   }       �// GetMarketingFeed is a generic endpoint that supports various marketing/out-of-app recommendation content with same content format, e.g. fetches fresh finds feed items rendering fresh finds page   �// Example design of content format: https://docs.google.com/document/d/1Ttmhht-bXOR7wc0-XA7IYBo52GNB_RQOy0vpmzhwhuk/edit#heading=h.nkalkkf3hlff   �func (fh FeedHandler) GetMarketingFeed(ctx thrift.Context, request *eatsfeedt.GetMarketingFeedRequest) (*eatsfeedt.GetMarketingFeedResponse, error) {   G	newCtx, span := jaegerhelpers.CreateChildSpan(ctx, "GetMarketingFeed")   	if span != nil {   		defer span.Finish()   	}       h	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertFeedParams(newCtx, request.FeedRequest)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   .		}).Error("Invalid GetMarketingFeed request")   		return nil, err   	}       &	if request.MarketingFeedType == nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"request": request,   G		}).Error("Invalid GetMarketingFeed request: Empty MarketingFeedType")   		return nil, nil   	}       &	if request.IsSetTargetingStoreTag() {   B		getFeedParams.TargetingStoreTag = request.GetTargetingStoreTag()   	}       "	if request.IsSetPromotionUuid() {   :		getFeedParams.PromotionUUID = request.GetPromotionUuid()   	}       "	if request.IsSetBillboardUuid() {   :		getFeedParams.BillboardUUID = request.GetBillboardUuid()   	}       k	marketingFeedResp, err := fh.FeedService.GetMarketingFeed(ctx, *getFeedParams, *request.MarketingFeedType)   	if err != nil {   U		fh.Deps.Logger().WithField("error", err.Error()).Error("Eatsfeed.GetMarketingFeed")   	}       	return marketingFeedResp, nil   }       Z//GetScheduledFeedDeliveryHours returns hour slots that are available for scheduled orders   �func (fh FeedHandler) GetScheduledFeedDeliveryHours(ctx thrift.Context, request *eatsfeedt.GetScheduledFeedDeliveryHoursRequest) (*eatsfeedt.GetScheduledFeedDeliveryHoursResponse, error) {   O	_, span := jaegerhelpers.CreateChildSpan(ctx, "GetScheduledFeedDeliveryHours")   	if span != nil {   		defer span.Finish()   	}   7	if request.GetTargetDeliveryLocationTimezone() == "" {   )		fh.Deps.Logger().WithFields(log.Fields{   			"request": request,   ]		}).Error("Invalid GetScheduledFeedDeliveryHours request: Empty delivery location timezone")   @		return &eatsfeedt.GetScheduledFeedDeliveryHoursResponse{}, nil   	}   {	deliveryHoursInfo := lib.GetDeliveryHoursInfos(clock.New(), request.GetTargetDeliveryLocationTimezone(), fh.Deps.Logger())   d	return &eatsfeedt.GetScheduledFeedDeliveryHoursResponse{DeliveryHoursInfos: deliveryHoursInfo}, nil   }       Ifunc (fh FeedHandler) getFeedTimeout(ctx context.Context) time.Duration {   `	timeout, _ := fh.Deps.Flipr().GetTimeDurationValue(ctx, getFeedTimeoutFlipr, flipr.Constraints{   ,		flipr.Constraints2(map[string]interface{}{   2			"environment": fh.Deps.Environment.Environment,   		}),   	}, getFeedTimeoutFallback)       	return timeout   }       Y//GetDeliveryCountdown populates restaurants in Delivery Countdown Hub with vertical view   W// https://docs.google.com/document/d/1sU8wtjEPXTKUmw8HjqRjaoe7PMIZJkgHN2AOkc7flKQ/edit   �func (fh FeedHandler) GetDeliveryCountdown(ctx thrift.Context, request *eatsfeedt.GetDeliveryCountDownRequest) (*eatsfeedt.GetDeliveryCountDownResponse, error) {   e	getFeedParams, err := lib.NewRequestConstructor(fh.Deps).ConvertFeedParams(ctx, request.FeedRequest)   	if err != nil {   )		fh.Deps.Logger().WithFields(log.Fields{   			"error":   err.Error(),   			"request": request,   9		}).Error("Invalid GetDeliveryCountDownRequest request")   		return nil, err   	}       D	return fh.FeedService.GetDeliveryCountDownFeed(ctx, *getFeedParams)   }5�5��