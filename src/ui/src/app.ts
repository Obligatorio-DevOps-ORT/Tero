import express from 'express';
import path from 'path';
import { createProxyMiddleware } from 'http-proxy-middleware';
import { log, requestLogger, registrarEventosDelProceso, proxyLogger } from './logger';

registrarEventosDelProceso();

const app = express();

// Primero el logger, antes de los proxies: así registra también los requests
// que los proxies responden sin pasar por el resto de la app.
app.use(requestLogger);

const CATALOG_URL  = process.env.RETAIL_UI_ENDPOINTS_CATALOG  || 'http://catalog:8080';
const CARTS_URL    = process.env.RETAIL_UI_ENDPOINTS_CARTS    || 'http://carts:8080';
const CHECKOUT_URL = process.env.RETAIL_UI_ENDPOINTS_CHECKOUT || 'http://checkout:8080';
const ORDERS_URL   = process.env.RETAIL_UI_ENDPOINTS_ORDERS   || 'http://orders:8080';

// Express strips the mount path before the proxy sees it.
// pathRewrite re-adds the service prefix, stripping trailing slash on bare paths.
const rewrite = (prefix: string) => (p: string) => prefix + (p === '/' ? '' : p);

app.use('/api/catalog',  createProxyMiddleware({ target: CATALOG_URL,  changeOrigin: true, pathRewrite: rewrite('/catalog'), logger: proxyLogger }));
app.use('/api/carts',    createProxyMiddleware({ target: CARTS_URL,    changeOrigin: true, pathRewrite: rewrite('/carts'), logger: proxyLogger }));
app.use('/api/checkout', createProxyMiddleware({ target: CHECKOUT_URL, changeOrigin: true, pathRewrite: rewrite('/checkout'), logger: proxyLogger }));
app.use('/api/orders',   createProxyMiddleware({ target: ORDERS_URL,   changeOrigin: true, pathRewrite: rewrite('/orders'), logger: proxyLogger }));

app.use(express.static(path.join(__dirname, '../public')));

app.get('/health', (_req, res) => res.send('OK'));

const port = parseInt(process.env.PORT || '8080', 10);
app.listen(port, () => log('INFO', 'servidor escuchando', { port }));
