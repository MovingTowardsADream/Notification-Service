# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [protos/proto/notify/notify.proto](#protos_proto_notify_notify-proto)
    - [Channels](#notify-Channels)
    - [MailApproval](#notify-MailApproval)
    - [MailNotify](#notify-MailNotify)
    - [PhoneApproval](#notify-PhoneApproval)
    - [PhoneNotify](#notify-PhoneNotify)
    - [Preferences](#notify-Preferences)
    - [SendMessageRequest](#notify-SendMessageRequest)
    - [SendMessageResponse](#notify-SendMessageResponse)
    - [UserPreferencesRequest](#notify-UserPreferencesRequest)
    - [UserPreferencesResponse](#notify-UserPreferencesResponse)
  
    - [NotifyType](#notify-NotifyType)
  
    - [SendNotify](#notify-SendNotify)
    - [Users](#notify-Users)
  
- [Scalar Value Types](#scalar-value-types)



<a name="protos_proto_notify_notify-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## protos/proto/notify/notify.proto



<a name="notify-Channels"></a>

### Channels



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| mail | [MailNotify](#notify-MailNotify) |  |  |
| phone | [PhoneNotify](#notify-PhoneNotify) |  |  |






<a name="notify-MailApproval"></a>

### MailApproval



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| approval | [bool](#bool) |  |  |






<a name="notify-MailNotify"></a>

### MailNotify



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subject | [string](#string) |  |  |
| body | [string](#string) |  |  |






<a name="notify-PhoneApproval"></a>

### PhoneApproval



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| approval | [bool](#bool) |  |  |






<a name="notify-PhoneNotify"></a>

### PhoneNotify



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| body | [string](#string) |  |  |






<a name="notify-Preferences"></a>

### Preferences



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| mail | [MailApproval](#notify-MailApproval) |  |  |
| phone | [PhoneApproval](#notify-PhoneApproval) |  |  |






<a name="notify-SendMessageRequest"></a>

### SendMessageRequest
The SendMessageRequest message contains the information necessary to send the notification


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| userId | [string](#string) |  |  |
| notifyType | [NotifyType](#notify-NotifyType) |  |  |
| channels | [Channels](#notify-Channels) |  |  |






<a name="notify-SendMessageResponse"></a>

### SendMessageResponse
The SendMessageResponse message contains information about the result of sending the notification


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| respond | [string](#string) |  |  |






<a name="notify-UserPreferencesRequest"></a>

### UserPreferencesRequest
The UserPreferencesRequest message contains information about the user&#39;s preferences


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| userId | [string](#string) |  |  |
| preferences | [Preferences](#notify-Preferences) |  |  |






<a name="notify-UserPreferencesResponse"></a>

### UserPreferencesResponse
The UserPreferencesResponse message contains information about the result of changing the user&#39;s preferences


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| respond | [string](#string) |  |  |





 


<a name="notify-NotifyType"></a>

### NotifyType


| Name | Number | Description |
| ---- | ------ | ----------- |
| moderate | 0 |  |
| significant | 1 |  |
| alert | 2 |  |


 

 


<a name="notify-SendNotify"></a>

### SendNotify
The SendNotify service provides methods for sending notifications

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SendMessage | [SendMessageRequest](#notify-SendMessageRequest) | [SendMessageResponse](#notify-SendMessageResponse) |  |


<a name="notify-Users"></a>

### Users
The Users service provides methods for working with user preferences

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| UserPreferences | [UserPreferencesRequest](#notify-UserPreferencesRequest) | [UserPreferencesResponse](#notify-UserPreferencesResponse) |  |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

