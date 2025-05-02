// assets/js/shopify.js

/**
 * Migrated to Shopify Storefront API (no JS Buy SDK)
 * Requires window.shopifyConfig with:
 *   - shopifyDomain: your-shop.myshopify.com
 *   - storefrontAccessToken: Storefront API token (read products/public)
 *
 * Usage: setupShopifyBuyButton(productId, productVariantId)
 */

function setupShopifyBuyButton(productId, productVariantId) {
    console.log('[Shopify] setupShopifyBuyButton called with:', { productId, variantId: productVariantId });
    const shopifyConfig = window.shopifyConfig;
    const endpoint = `https://${shopifyConfig.shopifyDomain}/api/2023-07/graphql.json`;
    const accessToken = shopifyConfig.storefrontAccessToken;

    let buyButton = document.getElementById('buy-button');
    let buyButtonSpinner = document.getElementById('buy-button-spinner');
    let unavailableToBuyButton = document.getElementById('unavailable-to-buy-button');
    let buyButtonProcess = document.getElementById('buy-button-process');

    // fixed page back from browser cache
    window.addEventListener("pageshow", function(event) {
        let historyTraversal = event.persisted ||
            (typeof window.performance !== "undefined" &&
                window.performance.navigation.type === 2);
        if (historyTraversal) {
            setupBuyButton();
        }
    });

    async function fetchVariantAvailability(productId, variantId) {
        console.log('[Shopify] Checking availability for:', { productId, variantId });
        const query = `
        query ($id: ID!) {
          product(id: $id) {
            variants(first: 50) {
              edges {
                node {
                  id
                  availableForSale
                }
              }
            }
          }
        }`;
        const variables = { id: productId };
        const res = await fetch(endpoint, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "X-Shopify-Storefront-Access-Token": accessToken,
            },
            body: JSON.stringify({ query, variables })
        });
        const data = await res.json();
        console.log('[Shopify] Storefront API response:', data);
        if (!data.data || !data.data.product) {
            console.warn('[Shopify] No product data found in response');
            return false;
        }
        const variants = data.data.product.variants.edges;
        console.log('[Shopify] Variants returned:', variants.map(v => v.node.id));
        const found = variants.find(v => v.node.id === variantId);
        if (!found) {
            console.warn('[Shopify] Variant not found for ID:', variantId);
        } else {
            console.log('[Shopify] Variant found:', found.node);
        }
        return found ? found.node.availableForSale : false;
    }

    async function setupBuyButton() {
        buyButton.classList.add('hidden');
        buyButtonProcess.classList.add('hidden');
        buyButtonSpinner.classList.remove('hidden');
        unavailableToBuyButton.classList.add('hidden');

        // Check variant availability via Storefront API
        let isAvailable = false;
        try {
            isAvailable = await fetchVariantAvailability(productId, productVariantId);
        } catch (e) {
            console.error('Error fetching product availability:', e);
        }

        buyButtonSpinner.classList.add('hidden');
        if (!isAvailable) {
            unavailableToBuyButton.classList.remove('hidden');
            return;
        }
        buyButton.classList.remove('hidden');
    }

    setupBuyButton();

    buyButton.addEventListener('click', function() {
        buyButton.classList.add('hidden');
        buyButtonProcess.classList.remove('hidden');
        try {
            // productVariantId is already in the format 'gid://shopify/ProductVariant/12345678'
            const matches = productVariantId.match(/ProductVariant\/(\d+)$/);
            if (!matches || matches.length < 2) {
                throw new Error('Could not parse variant ID: ' + productVariantId);
            }
            const numericVariantId = matches[1];
            // Construct the direct checkout URL
            const checkoutUrl = `https://${shopifyConfig.shopifyDomain}/cart/${numericVariantId}:1`;
            window.location.href = checkoutUrl;
        } catch (error) {
            console.error('Error processing checkout:', error);
            buyButtonProcess.classList.add('hidden');
            buyButton.classList.remove('hidden');
        }
        setTimeout(function() {
            buyButtonProcess.classList.add('hidden');
            buyButton.classList.remove('hidden');
            console.error('Checkout redirect timeout - please try again');
        }, 3000);
    });
}