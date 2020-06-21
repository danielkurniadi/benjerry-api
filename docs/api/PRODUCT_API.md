# Product API Schema

This documenation includes all API endpoints for CRUD Ben & Jerry products.

Notes:

- Timestamp used is in seconds. (i.e. need to times with 1000 in JavaScript)
- Every failed request is expected to be responded with an `error_message` field. For example:

```json
{
  "Message": "Insufficient permissions"
}
```
---

## Get Product Information

`GET api/products/<product_id>`

Permission Level: Read Permission, all member.

#### Response

##### No Error
`HTTP 200 OK`

| Name                  | Value                 | Description
| -----------------     | --------              | -----------
| `product`             | `Product Object`      | Details of the object

```json
{
  "product": {
      "productId": "646",
      "name": "Vanilla Toffee Bar Crunch",
      "image_closed": "/files/live/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
      "image_open": "/files/live/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
      "description": "Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
      "story": "Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars .... We love it and know you will too!",
      "sourcing_values": [
          "Non-GMO",
          "Cage-Free Eggs",
          "Fairtrade",
          "Responsibly Sourced Packaging",
          "Caring Dairy"
      ],
      "ingredients": [
          "cream",
          "skim milk",
          "liquid sugar",
          "water",
          "sugar",
          "coconut oil",
          "natural flavor",
          "salt",
          "vegetable oil"
       ],
       
      "allergy_info": "may contain wheat, peanuts and other tree nuts",
      "dietary_certifications": "Kosher",
   }
}
```

##### Error
`HTTP 404 Not Found`
| Name                  | Value                 | Description
| -----------------     | --------              | -----------
| `message`              | `False`              | Unsuccessful, due to resource not found


---

## Create Product Information

`POST api/products/`

Permission Level: Write Permission, admin only.

#### Request

*Header*: 
```
"session_token=02041fd4-c086-4e54-b931-be55a6043bd9"
```
*Body:*
```json
{
  "name": "Vanilla Toffee Bar Crunch",
  "image_closed": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing.png",
  "image_open": "/files/live/sites/systemsite/files/flavors/products/us/pint/open-closed-pints/vanilla-toffee-landing-open.png",
  "description": "Vanilla Ice Cream with Fudge-Covered Toffee Pieces",
  "story": "Vanilla What Bar Crunch? We gave this flavor a new name to go with the new toffee bars we’re using as part of our commitment to source Fairtrade Certified and non-GMO ingredients. We love it and know you will too!",
  "sourcing_values": [
    "Non-GMO",
    "Cage-Free Eggs",
    "Fairtrade",
    "Responsibly Sourced Packaging",
    "Caring Dairy"
  ],
  "ingredients": [
    "cream",
    "skim milk",
    "liquid sugar",
    "water",
    "sugar",
    "coconut oil",
    "egg yolks",
    "butter",
    "vanilla extract",
    "almonds",
    "cocoa (processed with alkali)",
    "milk",
    "soy lecithin",
    "cocoa",
    "natural flavor",
    "salt",
    "vegetable oil",
    "guar gum",
    "carrageenan"
  ],
  "allergy_info": "may contain wheat, peanuts and other tree nuts",
  "dietary_certifications": "Kosher",
  "productId": "646"
}
```

#### Response

##### No Error
`HTTP 201 CREATED`

| Name                  | Value                 | Description
| -----------------     | --------              | -----------
| `Message`             | `String`              | Successfully created


---

## Update Product Information

`PUT api/products/<product_id>`

Permission Level: Edit Permission, admin only.

#### Request

*Header*: 
```
"session_token=02041fd4-c086-4e54-b931-be55a6043bd9"
```
*Body:*
```json
{
  "description": "Edited: Mactha Vanilla Toffe Bar? Why not!",
  "story": "Edited: Vanilla What Bar Crunch? Add some matcha and toffe bar. We love it and know you will too!",
  "sourcing_values": [
    "Non-GMO",
    "Cage-Free Eggs"
  ]
}
```
> Note: The request body fields for `PUT` are similar to `POST` Create Body. However, the `PUT` are flexible, you only specify field(s) that you wanted to update. `productID` fields cannot be updated. 


#### Response

##### No Error
`HTTP 200 OK`

| Name                  | Value                 | Description
| -----------------     | --------              | -----------
| `Message`             | `String`              | Successfully updated

---

## Delete Product Information

`DELETE api/products/<product_id>`

Permission Level: Delete Permission, admin only.

#### Response

##### No Error
`HTTP 200 OK`

| Name                  | Value                 | Description
| -----------------     | --------              | -----------
| `Message`             | `String`              | Successfully updated
