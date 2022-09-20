# COREL ***(Correlation in Distributed Logs)***
## Usages
There are multiple places to use corel for establishment of related logs in microservice architecture. but 2 of them are most important.
One while making intranet service-to-service call and the other while pushing any message to kafka topic or rabbitmq exchange.
### For Api Call
```go=
import "gitlab.com/dotpe/mindbenders/corel"
....
req, err := http.NewRequest("GET",url,nil)
if err != nil {
	return res, err
}
corel.AttachCorelToHttpFromCtx(ctx, req)
```

`corel.AttachCorelToHttpFromCtx(ctx, req)` is important while adding corel information as header for http call.
### For Pushing to Topic
##### publisher side
```go
import "gitlab.com/dotpe/mindbenders/corel"
...
func getProduceTaskRequest(ctx context.Context, data interface{}, topic string) (*ProduceTaskRequest, error) {
	corelId, _ := corel.GetCorelationId(ctx)
	reqData := ProduceTaskRequestData{
		Data: data,
		Correlation: corel.EncodeCorel(corelId.Child()),
	}

	reqDataBytes, err := json.Marshal(reqData)
	if err != nil {
		return  nil, errors.New("unable to create produce task request")
	}

	taskRequest := ProduceTaskRequest{
		ReqBytes: reqDataBytes,
		Topic: "topic.x.y.z",
	}
	return &taskRequest, nil
}
```
##### consumer side
```go
import "gitlab.com/dotpe/mindbenders/corel"

...
var tmp struct {
	CoRelationId string  `json:"correlation"`
	Order producer.UserOrderKafkaMsg `json:"data"`
}
if  err := json.Unmarshal(raw, &tmp); err != nil {
	return err
}
oe.CoRelationId = corel.DecodeCorelationId(tmp.CoRelationId).Sibling()
oe.Order = &tmp.Order
return  nil

.............
ctx := corel.NewCorelCtxFromCorel(oe.CoRelationId)

```
