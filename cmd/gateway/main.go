// ... داخل main ...

// 1. Setup WS Hub (مدیریت کانکشن‌ها)
// hub := ws.NewHub()
// go hub.Run()

// 2. Setup Redis Subscriber
// sub := worker.NewSubscriber(rdb, hub)
// go sub.Start(context.Background())

// 3. Setup Router & Run
// ...