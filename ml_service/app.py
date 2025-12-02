from fastapi import FastAPI, File, UploadFile, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from ultralytics import YOLO
import cv2
import numpy as np
import uvicorn
from typing import List, Dict
import logging

# Настройка логирования
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = FastAPI(
    title="SafePass ML Service",
    description="YOLO PPE Detection Microservice",
    version="1.0.0"
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # В production укажи конкретные origins
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Загрузка модели
MODEL_PATH = "../ml_model/runs/detect/train4/weights/best.pt"
try:
    model = YOLO(MODEL_PATH)
    logger.info(f"✅ Model loaded successfully from {MODEL_PATH}")
except Exception as e:
    logger.error(f"❌ Failed to load model: {e}")
    raise

# Классы нарушений
VIOLATION_CLASSES = ["no_helmet", "no_goggle", "no_gloves", "no_boots"]


@app.get("/")
async def root():
    """Health check"""
    return {
        "service": "SafePass ML Service",
        "status": "running",
        "model": MODEL_PATH,
        "version": "1.0.0"
    }


@app.get("/health")
async def health():
    """Detailed health check"""
    return {
        "status": "healthy",
        "model_loaded": model is not None,
        "model_path": MODEL_PATH
    }


@app.post("/detect")
async def detect_ppe(file: UploadFile = File(...)):
    """
    Детекция СИЗ на загруженном изображении
    
    Args:
        file: Изображение (JPG, PNG)
    
    Returns:
        {
            "detections": [...],
            "has_violations": bool,
            "violations": [...]
        }
    """
    try:
        # Чтение изображения
        contents = await file.read()
        nparr = np.frombuffer(contents, np.uint8)
        img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
        
        if img is None:
            raise HTTPException(status_code=400, detail="Invalid image file")
        
        # YOLO детекция
        logger.info(f"Running detection on image: {file.filename}")
        results = model.predict(img, conf=0.35, verbose=False)
        
        # Парсинг результатов
        detections = []
        for result in results:
            boxes = result.boxes
            for box in boxes:
                cls = int(box.cls[0])
                conf = float(box.conf[0])
                label = model.names[cls]
                bbox = box.xyxy[0].tolist()
                
                detections.append({
                    "class": label,
                    "confidence": round(conf, 3),
                    "bbox": [round(x, 2) for x in bbox]
                })
        
        # Фильтрация нарушений
        violations = [d for d in detections if d["class"] in VIOLATION_CLASSES]
        
        logger.info(f"✅ Detection complete: {len(detections)} objects, {len(violations)} violations")
        
        return {
            "detections": detections,
            "has_violations": len(violations) > 0,
            "violations": violations,
            "total_detections": len(detections),
            "image_size": {"width": img.shape[1], "height": img.shape[0]}
        }
    
    except Exception as e:
        logger.error(f"❌ Detection error: {e}")
        raise HTTPException(status_code=500, detail=f"Detection failed: {str(e)}")


@app.post("/detect-batch")
async def detect_batch(files: List[UploadFile] = File(...)):
    """
    Детекция на нескольких изображениях
    """
    results = []
    
    for file in files:
        try:
            contents = await file.read()
            nparr = np.frombuffer(contents, np.uint8)
            img = cv2.imdecode(nparr, cv2.IMREAD_COLOR)
            
            if img is None:
                results.append({
                    "filename": file.filename,
                    "error": "Invalid image"
                })
                continue
            
            detections = []
            yolo_results = model.predict(img, conf=0.35, verbose=False)
            
            for result in yolo_results:
                for box in result.boxes:
                    detections.append({
                        "class": model.names[int(box.cls[0])],
                        "confidence": round(float(box.conf[0]), 3)
                    })
            
            violations = [d for d in detections if d["class"] in VIOLATION_CLASSES]
            
            results.append({
                "filename": file.filename,
                "detections": len(detections),
                "violations": len(violations),
                "has_violations": len(violations) > 0
            })
        
        except Exception as e:
            results.append({
                "filename": file.filename,
                "error": str(e)
            })
    
    return {"results": results, "total": len(files)}


if __name__ == "__main__":
    logger.info("🚀 Starting SafePass ML Service...")
    uvicorn.run(
        app,
        host="0.0.0.0",
        port=8001,
        log_level="info"
    )