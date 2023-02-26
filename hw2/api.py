import typing
import base64
from fastapi import FastAPI, HTTPException, UploadFile
from fastapi.responses import Response
from pydantic import BaseModel
from io import BytesIO


class Product(BaseModel):
    id: typing.Union[str, None]
    name: str
    description: typing.Union[str, None]


app = FastAPI()
products: typing.Dict[int, Product] = {}
images: typing.Dict[int, bytes] = {}
free_id = 0


@app.get("/products")
def get_products():
    return list(products.values())


@app.get("/products/{product_id}")
def get_product_item(product_id: int):
    if product_id not in products:
        raise HTTPException(status_code=404, detail="Product not found")
    return products[product_id]


@app.post("/products", status_code=201)
def add_product_item(item: Product):
    global free_id
    item.id = free_id
    free_id += 1
    products[item.id] = item
    return item


@app.post("/products/{product_id}/image", status_code=201)
async def upload_product_image(product_id: int, image: UploadFile):
    if product_id not in products:
        raise HTTPException(status_code=404, detail="Product not found")
    data = await image.read()
    images[product_id] = base64.b64encode(data)
    return {"status": "success"}


@app.get("/products/{product_id}/image", status_code=200, response_class=Response)
def load_product_image(product_id: int):
    if product_id not in products:
        raise HTTPException(status_code=404, detail="Product not found")
    if product_id not in images:
        raise HTTPException(status_code=404, detail="Product hasn't image")
    data = base64.b64decode(images[product_id])
    return Response(content=data, media_type="image/png")


@ app.put("/products/{product_id}", status_code=201)
def update_product_item(product_id: int, item: Product):
    if product_id not in products:
        raise HTTPException(status_code=404, detail="Product not found")

    item.id = product_id
    products[product_id] = item
    return item


@ app.delete("/products/{product_id}", status_code=201)
def remove_product_item(product_id: int):
    if product_id not in products:
        raise HTTPException(status_code=404, detail="Product not found")

    del products[product_id]
    return {"status": "success"}
