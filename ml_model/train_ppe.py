from ultralytics import YOLO

def main():
    # Загружаем предобученную модель (можно заменить на yolo11n.pt)
    model = YOLO("yolov8n.pt")

    # Обучаем модель
    model.train(
        data="data.yaml",    # путь к yaml
        epochs=50,           # количество эпох
        imgsz=640,           # размер картинок
        batch=8,             # размер batch (уменьши если мало VRAM)
        device="cuda"        # можно заменить на "cpu" если CUDA не работает
    )

if __name__ == "__main__":
    main()
