const axios = require('axios');
require('dotenv').config();

// Your Etsy API key
const api_key = process.env.ETSY_API_KEY;
const shop_id = process.env.SHOP_ID;
const shop_name = process.env.SHOP_NAME;

// The endpoint URL
const url = 'https://openapi.etsy.com/v2/listings/active';

const requestOptions = {
    'method': 'GET',
    'headers': {
        'x-api-key': api_key,
    },
};
async function pingEtsy() {
    const fetch = await import('node-fetch');
    const response = await fetch.default(
        'https://api.etsy.com/v3/application/openapi-ping',
        requestOptions
    );

    if (response.ok) {
        const data = await response.json();
        console.log(data);
    } else {
        const data = await response.json();
        console.log(data);
    }
}
async function getShopId() {
    const url = 'https://openapi.etsy.com/v3/application/shops?shop_name=' + shop_name;
    const requestOptions = {
        'method': 'GET',
        'headers': {
            'x-api-key': api_key,
        }
    };
    const fetch = await import('node-fetch');
    const response = await fetch.default(url, requestOptions);

    if (response.ok) {
        const data = await response.json();
        console.log(data);
    } else {
        const data = await response.json();
        console.log(data);
    }
}
async function findAllActiveListingsByShop() {
    const url = `https://openapi.etsy.com/v3/application/shops/${shop_id}/listings/active`;
    const requestOptions = {
        'method': 'GET',
        'headers': {
            'x-api-key': api_key,
        }
    };
    const fetch = await import('node-fetch');
    const response = await fetch.default(url, requestOptions);

    if (response.ok) {
        const data = await response.json();
        console.log(data);
    } else {
        const data = await response.json();
        console.log(data);
    }
}

// pingEtsy();
// getShopId();
findAllActiveListingsByShop();