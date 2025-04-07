import express from 'express';

const app = express();

app.get('/',(req, res) =>{
    res.send('Hello NodeJS!!');
});

const PORT = 3000;

app.listen(PORT, () => {
    console.log(`Servidor ejecutando en http://localhost:${PORT}`);
});