import React from 'react'

const Timelive = () => {
  return (
  
      <section className="mb-8">
        <h3 className="text-xl font-medium mb-3">Онлайн-трансляция с камер</h3>
        <div className="bg-black rounded-2xl overflow-hidden shadow-md h-72 flex items-center justify-center text-gray-400">
          <p>RTSP поток камеры (пример: rtsp://192.168.1.10/live)</p>
        </div>
      </section>
  )
}

export default Timelive