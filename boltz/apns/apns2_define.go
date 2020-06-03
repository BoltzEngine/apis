package apns

import (
	"net/http"
)

/*
Notice.
  Apple has announced that it will discontinue the binary interface from APNs from 2020-11.
  https://developer.apple.com/news/?id=11042019a

  The information in this file is subject to substantial revision in the future.
  Never refer to a private attribute variable or function from outside.
  There is no guarantee of compatibility.
*/

func isInvalidToken(status int, reason string) (ret bool) {
	return status == http.StatusBadRequest && reason == http2ReasonBadDeviceToken.String() ||
		status == http.StatusNotFound ||
		status == http.StatusGone
}

type apns2Reason string
type apns2ErrorDetail map[apns2Reason]string

func (r apns2Reason) String() (ret string) {
	return string(r)
}

func (r apns2Reason) Detail() (ret string) {
	return apns2ErrorDetails[r]
}

const (
	http2ReasonBadCollapseId               = apns2Reason("BadCollapseId")
	http2ReasonBadDeviceToken              = apns2Reason("BadDeviceToken")
	http2ReasonBadExpirationDate           = apns2Reason("BadExpirationDate")
	http2ReasonBadMessageId                = apns2Reason("BadMessageId")
	http2ReasonBadPriority                 = apns2Reason("BadPriority")
	http2ReasonBadTopic                    = apns2Reason("BadTopic")
	http2ReasonDeviceTokenNotForTopic      = apns2Reason("DeviceTokenNotForTopic")
	http2ReasonDuplicateHeaders            = apns2Reason("DuplicateHeaders")
	http2ReasonIdleTimeout                 = apns2Reason("IdleTimeout")
	http2ReasonMissingDeviceToken          = apns2Reason("MissingDeviceToken")
	http2ReasonMissingTopic                = apns2Reason("MissingTopic")
	http2ReasonPayloadEmpty                = apns2Reason("PayloadEmpty")
	http2ReasonTopicDisallowed             = apns2Reason("TopicDisallowed")
	http2ReasonBadCertificate              = apns2Reason("BadCertificate")
	http2ReasonBadCertificateEnvironment   = apns2Reason("BadCertificateEnvironment")
	http2ReasonExpiredProviderToken        = apns2Reason("ExpiredProviderToken")
	http2ReasonForbidden                   = apns2Reason("Forbidden")
	http2ReasonInvalidProviderToken        = apns2Reason("InvalidProviderToken")
	http2ReasonMissingProviderToken        = apns2Reason("MissingProviderToken")
	http2ReasonBadPath                     = apns2Reason("BadPath")
	http2ReasonMethodNotAllowed            = apns2Reason("MethodNotAllowed")
	http2ReasonUnregistered                = apns2Reason("Unregistered")
	http2ReasonPayloadTooLarge             = apns2Reason("PayloadTooLarge")
	http2ReasonTooManyProviderTokenUpdates = apns2Reason("TooManyProviderTokenUpdates")
	http2ReasonTooManyRequests             = apns2Reason("TooManyRequests")
	http2ReasonInternalServerError         = apns2Reason("InternalServerError")
	http2ReasonServiceUnavailable          = apns2Reason("ServiceUnavailable")
	http2ReasonShutdown                    = apns2Reason("Shutdown")
)

var apns2ErrorDetails = apns2ErrorDetail{
	http2ReasonBadCollapseId:               "The collapse identifier exceeds the maximum allowed size",
	http2ReasonBadDeviceToken:              "The specified device token was bad. Verify that the request contains a valid token and that the token matches the environment.",
	http2ReasonBadExpirationDate:           "The apns-expiration value is bad.",
	http2ReasonBadMessageId:                "The apns-id value is bad.",
	http2ReasonBadPriority:                 "The apns-priority value is bad.",
	http2ReasonBadTopic:                    "The apns-topic was invalid.",
	http2ReasonDeviceTokenNotForTopic:      "The device token does not match the specified topic.",
	http2ReasonDuplicateHeaders:            "One or more headers were repeated.",
	http2ReasonIdleTimeout:                 "Idle time out.",
	http2ReasonMissingDeviceToken:          "The device token is not specified in the request :path. Verify that the :path header contains the device token.",
	http2ReasonMissingTopic:                "The apns-topic header of the request was not specified and was required. The apns-topic header is mandatory when the client is connected using a certificate that supports multiple topics.",
	http2ReasonPayloadEmpty:                "The message payload was empty.",
	http2ReasonTopicDisallowed:             "Pushing to this topic is not allowed.",
	http2ReasonBadCertificate:              "The certificate was bad.",
	http2ReasonBadCertificateEnvironment:   "The client certificate was for the wrong environment.",
	http2ReasonExpiredProviderToken:        "The provider token is stale and a new token should be generated.",
	http2ReasonForbidden:                   "The specified action is not allowed.",
	http2ReasonInvalidProviderToken:        "The provider token is not valid or the token signature could not be verified.",
	http2ReasonMissingProviderToken:        "No provider certificate was used to connect to http and Authorization header was missing or no provider token was specified.",
	http2ReasonBadPath:                     "The request contained a bad :path value.",
	http2ReasonMethodNotAllowed:            "The specified :method was not POST.",
	http2ReasonUnregistered:                "The device token is inactive for the specified topic.",
	http2ReasonPayloadTooLarge:             "The message payload was too large. For regular remote notifications, the maximum size is 4KB. For VoIP notifications, the maximum size is 5KB.",
	http2ReasonTooManyProviderTokenUpdates: "The provider token is being updated too often.",
	http2ReasonTooManyRequests:             "Too many requests were made consecutively to the same device token.",
	http2ReasonInternalServerError:         "An internal server error occurred.",
	http2ReasonServiceUnavailable:          "The service is unavailable.",
	http2ReasonShutdown:                    "The server is shutting down.",
}
