// Convert string to snake_case
export function toSnakeCase(str) {
    return str
        .normalize('NFD').replace(/[\u0300-\u036f]/g, '')
        .replace(/[^a-zA-Z0-9_]+/g, ' ')
        .trim()
        .toLowerCase()
        .replace(/\s+/g, '_');
}

// Update retry policy preview
export function updateRetryPolicyPreview(type, maxAttempts, interval) {
    let previewText = '';
    
    if (type === 'constant') {
        previewText = `После каждой неудачи ожидать <strong>${interval}мс</strong> перед повтором.`;
    } else if (type === 'linear') {
        previewText = `После каждой неудачи ожидать <strong>${interval}мс × номер попытки</strong> перед повтором.`;
    } else if (type === 'exponential') {
        previewText = `После каждой неудачи ожидать <strong>${interval}мс × 2^(номер попытки-1)</strong> перед повтором.`;
    }
    
    if (maxAttempts) {
        previewText += ` Будет предпринято до <strong>${maxAttempts}</strong> попыток, затем задача будет считаться неуспешной.`;
    } else {
        previewText += ` Будет повторяться до успешного выполнения.`;
    }
    
    return previewText;
}

// Update pagination information
export function updatePaginationInfo(currentPage, itemsPerPage, total) {
    const start = (currentPage - 1) * itemsPerPage + 1;
    const end = Math.min(currentPage * itemsPerPage, total);
    
    return {
        start,
        end,
        total,
        hasPrev: currentPage > 1,
        hasNext: currentPage * itemsPerPage < total
    };
}

// Load HTML component
export async function loadComponent(elementId, componentPath) {
    try {
        const response = await fetch(componentPath);
        if (!response.ok) {
            throw new Error(`Failed to load component: ${response.status} ${response.statusText}`);
        }
        const html = await response.text();
        const element = document.getElementById(elementId);
        if (!element) {
            throw new Error(`Element with id "${elementId}" not found`);
        }
        element.innerHTML = html;
    } catch (error) {
        console.error('Error loading component:', error);
        const element = document.getElementById(elementId);
        if (element) {
            element.innerHTML = `<div class="text-red-500 p-4">Error loading component: ${error.message}</div>`;
        }
    }
} 