from ultralytics import YOLO

def main():
    # Загружаем модель
    model = YOLO("runs/detect/train4/weights/best.pt")

    # Запускаем детекцию в реальном времени с камеры
    model.predict(
        source=0,        # 0 = веб-камера
        show=True,       # показывать окно
        conf=0.35,        # минимальная уверенность
        save=False,      # не сохранять каждое изображение
        device='cuda'    # использовать GPU
    )

if __name__ == "__main__":
    main()
