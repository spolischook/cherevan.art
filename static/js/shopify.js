// assets/js/shopify.js

function setupShopifyBuyButton(productId, productVariantId, isInStock) {
    const shopifyConfig = window.shopifyConfig;

    const client = ShopifyBuy.buildClient({
        domain: shopifyConfig.shopifyDomain,
        storefrontAccessToken: shopifyConfig.storefrontAccessToken
    });

    let buyButton = document.getElementById('buy-button');
    let buyButtonSpinner = document.getElementById('buy-button-spinner');
    let unavailableToBuyButton = document.getElementById('unavailable-to-buy-button');
    let buyButtonProcess = document.getElementById('buy-button-process');

    // fixed page back from browser cache
    window.addEventListener("pageshow", function(event) {
        let historyTraversal = event.persisted ||
            (typeof window.performance != "undefined" &&
                window.performance.navigation.type === 2);
        if (historyTraversal) {
            setupBuyButton(isInStock);
        }
    });

    function setupBuyButton(isInStock) {
        buyButton.classList.add('hidden');
        buyButtonProcess.classList.add('hidden');
        buyButtonSpinner.classList.remove('hidden');
        unavailableToBuyButton.classList.add('hidden');
        if (!isInStock) {
            buyButtonSpinner.classList.add('hidden');
            unavailableToBuyButton.classList.remove('hidden');
            return;
        }

        client.product.fetch(productId).then((product) => {
            buyButtonSpinner.classList.add('hidden');

            if (!product || !product.variants || !product.variants[0].available) {
                unavailableToBuyButton.classList.remove('hidden');
                return;
            }

            buyButton.classList.remove('hidden');
        });
    }

    setupBuyButton(isInStock);

    buyButton.addEventListener('click', function() {
        buyButton.classList.add('hidden');
        buyButtonProcess.classList.remove('hidden');
        client.checkout.create().then((checkout) => {
            const lineItemsToAdd = [
                {
                    variantId: productVariantId,
                    quantity: 1,
                }
            ];

            client.checkout.addLineItems(checkout.id, lineItemsToAdd).then((checkout) => {
                window.location.href = checkout.webUrl;
            }).catch((error) => {
                console.error(error);
            });
        });
    });
}