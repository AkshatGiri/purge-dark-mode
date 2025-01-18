function Test() {
  return (
    <html lang="en">
      <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>Document</title>
      </head>
      <body>
        <div className="bg-white dark:bg-black">
          <h1 className="dark:text-white text-black">Hello World</h1>
          <div className={`bg-black dark:bg-black`}></div>
        </div>
      </body>
    </html>
  );
}
