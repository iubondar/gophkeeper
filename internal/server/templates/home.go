package templates

// HomePageHTML возвращает HTML код главной страницы
func HomePageHTML() string {
	return `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>GophKeeper API</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            margin: 0;
            padding: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container {
            background: white;
            border-radius: 12px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            padding: 40px;
            max-width: 800px;
            width: 90%;
            margin: 20px;
        }
        h1 {
            color: #333;
            text-align: center;
            margin-bottom: 30px;
            font-size: 2.5em;
            font-weight: 300;
        }
        .description {
            color: #666;
            text-align: center;
            margin-bottom: 40px;
            font-size: 1.1em;
        }
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }
        .feature {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        .feature h3 {
            color: #333;
            margin-top: 0;
            margin-bottom: 10px;
        }
        .feature p {
            color: #666;
            margin: 0;
        }
        .links {
            display: flex;
            justify-content: center;
            gap: 20px;
            flex-wrap: wrap;
        }
        .btn {
            display: inline-block;
            padding: 12px 24px;
            background: #667eea;
            color: white;
            text-decoration: none;
            border-radius: 6px;
            font-weight: 500;
            transition: all 0.3s ease;
            border: none;
            cursor: pointer;
        }
        .btn:hover {
            background: #5a6fd8;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }
        .btn-secondary {
            background: #6c757d;
        }
        .btn-secondary:hover {
            background: #5a6268;
        }
        .status {
            text-align: center;
            margin-top: 30px;
            padding: 15px;
            background: #d4edda;
            color: #155724;
            border-radius: 6px;
            border: 1px solid #c3e6cb;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🔐 GophKeeper API</h1>
        <p class="description">
            Безопасная клиент-серверная система для хранения конфиденциальной информации
        </p>
        
        <div class="features">
            <div class="feature">
                <h3>🔑 Аутентификация</h3>
                <p>JWT токены с refresh механизмом для безопасного доступа к данным</p>
            </div>
            <div class="feature">
                <h3>🔒 Шифрование</h3>
                <p>AES-256 шифрование на клиенте, передача только зашифрованных данных</p>
            </div>
            <div class="feature">
                <h3>📁 Типы данных</h3>
                <p>Логины/пароли, текстовые секреты, банковские карты, файлы</p>
            </div>
            <div class="feature">
                <h3>🔄 Синхронизация</h3>
                <p>Оптимистичные блокировки для разрешения конфликтов при одновременном доступе</p>
            </div>
        </div>
        
        <div class="links">
            <a href="/swagger-ui/index.html" class="btn">📖 Swagger UI</a>
            <a href="/swagger/swagger.json" class="btn btn-secondary">📄 JSON Schema</a>
        </div>
        
        <div class="status">
            ✅ Сервис работает нормально
        </div>
    </div>
</body>
</html>`
}
