free_id = 0


class Product:
    def __init__(self, name: str, description: str):
        self.id = free_id
        self.name = name
        self.description = description
        free_id += 1
