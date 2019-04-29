Carlos García de Marina Vilar - garciademarina@gmail.com

1. Please write a new Microservice called “Order Service” which will:
    a. Support the Standard Order Flow for an Ecommerce Order

    b. Consider an Order Object to have the following details 

        i. Order Id,
        ii. Total amount
        iii. Order Lines, where each line is a collection of
            1. Item SKU
            2. Price
            3. Quantity
        iv. Order Shipping Address
        v. Order Billing Addresss
    c. Consider the following to be standard order statuses

        i. Pending Confirmation
        ii. Confirmed
        iii. Sent to Warehouse
        iv. Shipped
        v. In Transit
        vi. Delivered
    d. The Microservice will provide relevant API calls for Order Creation, Status Updates
    e. The Microservice should generate relevant Events for each order status updates
    f. The service should have its own NoSQL persistence Layer


## Installation

### Using docker
```bash
cd $GOPATH/src/github.com/garciademarina/deporvillage
docker-compose build
docker-compose up -d 
```

### From source

Requires Go 1.12 or later.
Requires mongodb, check cmd/config.json in order to update mongodb url, port and database.

```bash
go get -u github.com/garciademarina/deporvillage
```

## Examples

Add new Order
```bash
curl -v -X POST "http://localhost:8080/order" -H "accept: application/json;" -d '{"id":2,"amount":10,"status":"cli","order_lines":[{"sku":"TROP-UA-PLAS-09","price":10,"quantity":1},{"sku":"TROP-NP-PLAS-65","price":10,"quantity":2},{"sku":"TROP-LT-PLAS-89","price":5,"quantity":10}],"shipping_address":{"first_name":"Jhon","last_name":"Snow","email":"j.snow@example.com","company":"Acme","phone":"555000111","line1":"711-2880 Nulla St.","line2":"","line3":"","city":"Mankato","country":"Mississippi","zip":"96522"},"billing_address":{"first_name":"Jhon","last_name":"Snow","email":"j.snow@example.com","company":"Acme","phone":"555000111","line1":"711-2880 Nulla St.","line2":"","line3":"","city":"Mankato","country":"Mississippi","zip":"96522"}}'
```

Update status Order
```bash
curl -v -X PUT "http://localhost:8080/order" -H "accept: application/json;" -d '{"id":2,"status":"confirmed"}'
```

Get Order by ID
```bash
curl -v -X GET "http://localhost:8080/order/2" -H "accept: application/json;"
```


## Api response
- Successfull request will return http status code 200 along with a json with additional information (or not).
- Failed request will return http status code 40X along with a json with additional information of the error.
```
Type (string) Posible values api_error
Code (string) Optional code
Message (string) A string representation for the error
```

## Structure
#### Server: 

Http server, set up all API routes using injected services from cmd/api/main.go

pkg/server

#### Services: 

/pkg/adding
/pkg/updating
/pkg/listing

#### Repository:

Repositories satisfy the *Repository* interface, this way we can swap mongodb by another noSQL database or SQL o in-memory implementation. 

/pkg/storage/mongodb

#### Handlers: 

http handler used by server.go, handlers will use injected services  from pkg/server.go

## Notes

- ShippingAddress and BillingAddress has their own struct, this way in the future we can get reports like: 
All order with "shipped" status where shipping city is Madrid.

- I asume "PUT /order" endpoint only updated the order status. If all order fields needs to be updated, it should be easy to implement by just adding the missing order fields.  

- Status must be one of the following: 
**pending_confirmation, confirmed, sent_to_warehouse, shipped, in_transit, delivered**

## endpoints 

#### GET /order/{ID}
Get the current order given an order ID. {ID}

**Arguments**
- ID (required) A string ID of an order

**Response**
Order object

```json
{
   "id":2,
   "amount":10,
   "status":"confirmed",
   "order_lines":[
      {
         "sku":"TROP-UA-PLAS-09",
         "price":10,
         "quantity":1
      },
      {
         "sku":"TROP-NP-PLAS-65",
         "price":10,
         "quantity":2
      },
      {
         "sku":"TROP-LT-PLAS-89",
         "price":5,
         "quantity":10
      }
   ],
   "shipping_address":{
      "first_name":"Jhon",
      "last_name":"Snow",
      "email":"j.snow@example.com",
      "company":"Acme",
      "phone":"555000111",
      "line1":"711-2880 Nulla St.",
      "line2":"",
      "line3":"",
      "city":"Mankato",
      "country":"Mississippi",
      "zip":"96522"
   },
   "billing_address":{
      "first_name":"Jhon",
      "last_name":"Snow",
      "email":"j.snow@example.com",
      "company":"Acme",
      "phone":"555000111",
      "line1":"711-2880 Nulla St.",
      "line2":"",
      "line3":"",
      "city":"Mankato",
      "country":"Mississippi",
      "zip":"96522"
   }
}
```

#### POST /order
Create a new order

**Arguments (json body)**
Order object

```json
{
   "id":2,
   "amount":10,
   "status":"confirmed",
   "order_lines":[
      {
         "sku":"TROP-UA-PLAS-09",
         "price":10,
         "quantity":1
      },
      {
         "sku":"TROP-NP-PLAS-65",
         "price":10,
         "quantity":2
      },
      {
         "sku":"TROP-LT-PLAS-89",
         "price":5,
         "quantity":10
      }
   ],
   "shipping_address":{
      "first_name":"Jhon",
      "last_name":"Snow",
      "email":"j.snow@example.com",
      "company":"Acme",
      "phone":"555000111",
      "line1":"711-2880 Nulla St.",
      "line2":"",
      "line3":"",
      "city":"Mankato",
      "country":"Mississippi",
      "zip":"96522"
   },
   "billing_address":{
      "first_name":"Jhon",
      "last_name":"Snow",
      "email":"j.snow@example.com",
      "company":"Acme",
      "phone":"555000111",
      "line1":"711-2880 Nulla St.",
      "line2":"",
      "line3":"",
      "city":"Mankato",
      "country":"Mississippi",
      "zip":"96522"
   }
}
```
**Response**

```json
{
    "success":"ok"
}
```

#### PUT /order
Update order status.

```json
{
    "id":2,
    "status":"confirmed"
}
```

**Response**

```json
{
    "success":"ok"
}
```

 


# Questions
* Once the order service is almost ready, how will you deploy the service to a cloud hosting? You can try to summarize your solution using 

- As I already have a dockerimage that contains only one executable file, it should be easy to deploy the container to kubernetes, I will need to upload the docker image to a repository and then just create the kubernetes job to deploy the uploaded image inside the kubernetes cluster. (I have limit knowlage of kubernetes)


a) CI/CD? How? Which tools?

Jenkings, (it's the only one I worked with). Using github webhooks, on every commit/merge to master branch should trigger the CI/CD pipeline. 
Once the webhook is received the docker image is updated, test are executed and if all test are successful then just update/updoad the docker image to the repository and deploy the new version to kubernetes.


b) Monitoring &amp; Alerts? Any concrete tool?

Prometheus, it's a monitoring system that also have an alerting system.

c) Logging?

ELK (elasticsearch , logstash, kibana)

- logstash, parse logs generated by the api.
- elasticsearch, used to store logs.
- kibana, used to explore and visualize logs. 

d) Deployment - Microplatform (AMI) vs Kubernetes?

With my limited experience with kubernetes, It shoud be easy to deploy de api to kubernetes.
Mongodb instance can be deployed to kubernetes or an external service selfmanaged. 
