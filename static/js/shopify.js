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
        
        try {
            // Decode the base64-encoded variant ID
            const decodedVariantId = atob(productVariantId);
            console.log('Decoded variant ID:', decodedVariantId);
            
            // Expected format: gid://shopify/ProductVariant/12345678
            // Extract just the numeric ID from the end
            const matches = decodedVariantId.match(/ProductVariant\/(\d+)$/);
            if (!matches || matches.length < 2) {
                throw new Error('Could not parse variant ID: ' + decodedVariantId);
            }
            
            const numericVariantId = matches[1];
            console.log('Numeric variant ID:', numericVariantId);
            
            // Construct the direct checkout URL
            const checkoutUrl = `https://${shopifyConfig.shopifyDomain}/cart/${numericVariantId}:1`;
            console.log('Checkout URL:', checkoutUrl);
            
            // Redirect to the checkout URL
            window.location.href = checkoutUrl;
            
        } catch (error) {
            console.error('Error processing checkout:', error);
            buyButtonProcess.classList.add('hidden');
            buyButton.classList.remove('hidden');
        }
        
        // Add error handling with a timeout in case the redirect fails
        setTimeout(function() {
            // If we're still here after 3 seconds, the redirect may have failed
            buyButtonProcess.classList.add('hidden');
            buyButton.classList.remove('hidden');
            console.error('Checkout redirect timeout - please try again');
        }, 3000);
    });
}