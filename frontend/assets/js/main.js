import { fetchExecutors, createExecutor, updateExecutor, getExecutor, getDLQTasks, clearDLQ } from './api.js';
import { toSnakeCase, updateRetryPolicyPreview, updatePaginationInfo } from './utils.js';

// Current state
let currentPage = 1;
const itemsPerPage = 20;
let currentExecutorId = null;
let currentStatusFilter = 'all';
let currentSearchQuery = '';
let isAddMode = false;
let currentDlqExecutor = null;

// DOM elements
let executorsTableBody;
let settingsDrawer;
let drawerPanel;
let closeDrawer;
let drawerBackdrop;
let cancelSettings;
let saveSettings;
let prevPage;
let nextPage;
let prevPageMobile;
let nextPageMobile;
let startItem;
let endItem;
let totalItems;
let dlqEnabled;
let dlqSettingsContainer;
let retryPolicyType;
let retryPolicyMaxAttempts;
let retryPolicyInterval;
let retryPolicyPreview;
let statusFilter;
let searchInput;
let dlqModal;
let closeDlqModal;
let dlqModalTitle;
let dlqDownloadBtn;
let dlqClearBtn;
let writeConcern;

// Initialize DOM elements
function initializeDOMElements() {
    executorsTableBody = document.getElementById('executorsTableBody');
    settingsDrawer = document.getElementById('settingsDrawer');
    drawerPanel = document.getElementById('drawerPanel');
    closeDrawer = document.getElementById('closeDrawer');
    drawerBackdrop = document.getElementById('drawerBackdrop');
    cancelSettings = document.getElementById('cancelSettings');
    saveSettings = document.getElementById('saveSettings');
    prevPage = document.getElementById('prevPage');
    nextPage = document.getElementById('nextPage');
    prevPageMobile = document.getElementById('prevPageMobile');
    nextPageMobile = document.getElementById('nextPageMobile');
    startItem = document.getElementById('startItem');
    endItem = document.getElementById('endItem');
    totalItems = document.getElementById('totalItems');
    dlqEnabled = document.getElementById('dlqEnabled');
    dlqSettingsContainer = document.getElementById('dlqSettingsContainer');
    retryPolicyType = document.getElementById('retryPolicyType');
    retryPolicyMaxAttempts = document.getElementById('retryPolicyMaxAttempts');
    retryPolicyInterval = document.getElementById('retryPolicyInterval');
    retryPolicyPreview = document.getElementById('retryPolicyPreview');
    statusFilter = document.getElementById('statusFilter');
    searchInput = document.getElementById('searchInput');
    dlqModal = document.getElementById('dlqModal');
    closeDlqModal = document.getElementById('closeDlqModal');
    dlqModalTitle = document.getElementById('dlqModalTitle');
    dlqDownloadBtn = document.getElementById('dlqDownloadBtn');
    dlqClearBtn = document.getElementById('dlqClearBtn');
    writeConcern = document.getElementById('writeConcern');
}

// Filter executors based on status and search query
function filterExecutors(executors) {
    return executors.filter(executor => {
        const statusMatch = currentStatusFilter === 'all' || 
                          (currentStatusFilter === 'enabled' && executor.enabled) ||
                          (currentStatusFilter === 'disabled' && !executor.enabled);
        
        const searchMatch = currentSearchQuery === '' || 
                          executor.name.toLowerCase().includes(currentSearchQuery);
        
        return statusMatch && searchMatch;
    });
}

// Render executors table
async function renderExecutors() {
    executorsTableBody.innerHTML = '';
    
    const executors = await fetchExecutors();
    const filteredExecutors = filterExecutors(executors);
    const totalFiltered = filteredExecutors.length;
    
    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const paginatedExecutors = filteredExecutors.slice(startIndex, endIndex);
    
    paginatedExecutors.forEach(executor => {
        const row = document.createElement('tr');
        row.className = 'soft-table-row table-row-hover';
        row.style.userSelect = 'none';
        row.dataset.id = executor.id;
        row.innerHTML = `
            <td class="px-6 py-4 whitespace-nowrap table-cell font-semibold text-gray-900 flex items-center gap-2">
                <i class="fas fa-pen text-blue-400 text-lg mr-2"></i>
                <span>${executor.name}</span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap table-cell text-gray-700">${executor.enabled 
                ? '<span class="px-2 inline-flex text-sm leading-5 font-semibold rounded-full bg-green-100 text-green-800">Включён</span>'
                : '<span class="px-2 inline-flex text-sm leading-5 font-semibold rounded-full bg-red-100 text-red-800">Отключён</span>'}
            </td>
            <td class="px-6 py-4 whitespace-nowrap table-cell text-center">
                <button title="Посмотреть логи"
                    class="log-btn flex items-center justify-center border border-blue-100 bg-white hover:bg-blue-100 text-blue-500 hover:text-blue-700 font-semibold w-32 h-10 rounded-xl shadow transition"
                    data-name="${executor.name}">
                    <i class="fas fa-file-alt text-xl"></i>
                </button>
            </td>
            <td class="px-6 py-4 whitespace-nowrap table-cell text-center">
                <button title="Посмотреть метрики"
                    class="metrics-btn flex items-center justify-center border border-blue-100 bg-white hover:bg-blue-100 text-blue-500 hover:text-blue-700 font-semibold w-32 h-10 rounded-xl shadow transition"
                    data-name="${executor.name}">
                    <i class="fas fa-chart-bar text-xl"></i>
                </button>
            </td>
            <td class="px-6 py-4 whitespace-nowrap table-cell text-center">
                ${
                    executor.config && executor.config.dlqConfig && executor.config.dlqConfig.enabled
                    ? `<button title="DLQ"
                        class="dlq-btn flex items-center justify-center border border-blue-100 bg-white hover:bg-blue-100 text-blue-500 hover:text-blue-700 font-semibold w-32 h-10 rounded-xl shadow transition"
                        data-name="${executor.name}">
                        <i class="fas fa-database text-xl"></i>
                    </button>`
                    : ''
                }
            </td>
        `;
        row.addEventListener('click', (e) => {
            if (e.target.closest('button')) return;
            currentExecutorId = row.dataset.id;
            openSettingsDrawer(currentExecutorId);
        });
        executorsTableBody.appendChild(row);
    });
    
    // Add event listeners for buttons
    document.querySelectorAll('.log-btn').forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            alert('Логи обработчика: ' + this.dataset.name);
        });
    });
    document.querySelectorAll('.metrics-btn').forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            alert('Метрики обработчика: ' + this.dataset.name);
        });
    });
    document.querySelectorAll('.dlq-btn').forEach(btn => {
        btn.addEventListener('click', function(e) {
            e.stopPropagation();
            currentDlqExecutor = this.dataset.name;
            dlqModalTitle.textContent = `DLQ для ${currentDlqExecutor}`;
            dlqModal.classList.remove('hidden');
        });
    });
    
    const paginationInfo = updatePaginationInfo(currentPage, itemsPerPage, filteredExecutors.length);
    startItem.textContent = paginationInfo.start;
    endItem.textContent = paginationInfo.end;
    totalItems.textContent = paginationInfo.total;
    
    prevPage.disabled = !paginationInfo.hasPrev;
    nextPage.disabled = !paginationInfo.hasNext;
    prevPageMobile.disabled = !paginationInfo.hasPrev;
    nextPageMobile.disabled = !paginationInfo.hasNext;
}

// Open settings drawer
async function openSettingsDrawer(executorId) {
    isAddMode = false;
    try {
        console.log('Opening settings drawer for executor:', executorId);
        const executor = await getExecutor(executorId);
        if (!executor) {
            console.error('No executor data received');
            return;
        }
        
        console.log('Opening settings for executor:', executor);
        
        currentExecutorId = executorId;
        const nameInput = document.getElementById('executorName');
        nameInput.value = executor.name;
        nameInput.readOnly = true;
        nameInput.classList.add('bg-gray-100', 'cursor-not-allowed');
        document.getElementById('executorNameNote').classList.remove('hidden');
        document.getElementById('enabled').checked = executor.enabled;
        
        const config = executor.config;
        console.log('Executor config:', config);
        
        if (!config) {
            console.error('No config found in executor data');
            throw new Error('No config found in executor data');
        }
        
        if (!config.retryPolicy || !config.dlqConfig || !config.writeConcern) {
            console.error('Missing required config fields:', config);
            throw new Error('Missing required config fields');
        }
        
        const retryPolicyTypeValue = config.retryPolicy.type.toLowerCase();
        console.log('Setting retry policy type:', retryPolicyTypeValue);
        retryPolicyType.value = retryPolicyTypeValue;
        
        retryPolicyMaxAttempts.value = config.retryPolicy.maxAttempts || '';
        retryPolicyInterval.value = parseInt(config.retryPolicy.interval.replace('s', '000')) || 1000;
        retryPolicyPreview.innerHTML = updateRetryPolicyPreview(
            retryPolicyTypeValue,
            config.retryPolicy.maxAttempts,
            parseInt(config.retryPolicy.interval.replace('s', '000'))
        );
        
        dlqEnabled.checked = config.dlqConfig.enabled;
        toggleDlqSettings(config.dlqConfig.enabled);
        document.getElementById('dlqQueueName').value = config.dlqConfig.queueName || '';
        
        const writeConcernLevel = config.writeConcern.level.toLowerCase();
        console.log('Setting write concern level:', writeConcernLevel);
        writeConcern.value = writeConcernLevel;
        
        document.getElementById('modalTitle').textContent = `Настройки для ${executor.name}`;
        settingsDrawer.classList.remove('hidden');
        setTimeout(() => drawerPanel.classList.add('open'), 10);
    } catch (error) {
        console.error('Error opening settings:', error);
        alert('Failed to load executor settings: ' + error.message);
    }
}

// Close settings drawer
function closeSettingsDrawer() {
    drawerPanel.classList.remove('open');
    setTimeout(() => settingsDrawer.classList.add('hidden'), 300);
}

// Save executor settings
async function saveExecutorSettings() {
    try {
        const config = {
            name: document.getElementById('executorName').value,
            enabled: document.getElementById('enabled').checked,
            retryPolicy: {
                type: retryPolicyType.value.toUpperCase(),
                maxAttempts: retryPolicyMaxAttempts.value ? parseInt(retryPolicyMaxAttempts.value) : 0,
                interval: retryPolicyInterval.value + 's'
            },
            dlqConfig: {
                enabled: dlqEnabled.checked,
                queueName: dlqEnabled.checked ? document.getElementById('dlqQueueName').value : ''
            },
            writeConcern: {
                level: document.getElementById('writeConcern').value.toUpperCase()
            }
        };

        if (isAddMode) {
            await createExecutor(config);
        } else {
            await updateExecutor(currentExecutorId, config);
        }

        alert('Settings saved successfully');
        closeSettingsDrawer();
        await renderExecutors();
    } catch (error) {
        console.error('Error saving settings:', error);
        alert('Failed to save executor settings: ' + error.message);
    }
}

// Navigate to previous page
function goToPrevPage() {
    if (currentPage > 1) {
        currentPage--;
        renderExecutors();
    }
}

// Navigate to next page
function goToNextPage() {
    const filteredExecutors = filterExecutors(executors);
    if (currentPage * itemsPerPage < filteredExecutors.length) {
        currentPage++;
        renderExecutors();
    }
}

// Toggle DLQ settings visibility
function toggleDlqSettings(enabled) {
    if (enabled) {
        dlqSettingsContainer.classList.remove('hidden');
    } else {
        dlqSettingsContainer.classList.add('hidden');
    }
}

// Set up event listeners
function setupEventListeners() {
    closeDrawer.addEventListener('click', closeSettingsDrawer);
    drawerBackdrop.addEventListener('click', closeSettingsDrawer);
    cancelSettings.addEventListener('click', closeSettingsDrawer);
    saveSettings.addEventListener('click', saveExecutorSettings);
    prevPage.addEventListener('click', goToPrevPage);
    nextPage.addEventListener('click', goToNextPage);
    prevPageMobile.addEventListener('click', goToPrevPage);
    nextPageMobile.addEventListener('click', goToNextPage);
    
    dlqEnabled.addEventListener('change', (e) => {
        toggleDlqSettings(e.target.checked);
    });
    
    retryPolicyType.addEventListener('change', () => {
        retryPolicyPreview.innerHTML = updateRetryPolicyPreview(
            retryPolicyType.value,
            retryPolicyMaxAttempts.value,
            parseInt(retryPolicyInterval.value)
        );
    });
    retryPolicyMaxAttempts.addEventListener('input', () => {
        retryPolicyPreview.innerHTML = updateRetryPolicyPreview(
            retryPolicyType.value,
            retryPolicyMaxAttempts.value,
            parseInt(retryPolicyInterval.value)
        );
    });
    retryPolicyInterval.addEventListener('input', () => {
        retryPolicyPreview.innerHTML = updateRetryPolicyPreview(
            retryPolicyType.value,
            retryPolicyMaxAttempts.value,
            parseInt(retryPolicyInterval.value)
        );
    });
    
    statusFilter.addEventListener('change', (e) => {
        currentStatusFilter = e.target.value;
        currentPage = 1;
        renderExecutors();
    });
    
    searchInput.addEventListener('input', (e) => {
        currentSearchQuery = e.target.value.toLowerCase();
        currentPage = 1;
        renderExecutors();
    });

    dlqModal.addEventListener('click', (e) => {
        if (e.target === dlqModal) {
            dlqModal.classList.add('hidden');
            currentDlqExecutor = null;
        }
    });

    dlqDownloadBtn.addEventListener('click', async () => {
        try {
            const tasks = await getDLQTasks(currentDlqExecutor);
            const blob = new Blob([JSON.stringify(tasks, null, 2)], { type: 'application/json' });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `${currentDlqExecutor}_dlq.json`;
            document.body.appendChild(a);
            a.click();
            window.URL.revokeObjectURL(url);
            document.body.removeChild(a);
        } catch (error) {
            console.error('Error downloading DLQ:', error);
            alert('Failed to download DLQ tasks');
        }
    });

    dlqClearBtn.addEventListener('click', async () => {
        if (!confirm('Are you sure you want to clear the DLQ?')) return;
        try {
            await clearDLQ(currentDlqExecutor);
            dlqModal.classList.add('hidden');
            currentDlqExecutor = null;
            alert('DLQ cleared successfully');
        } catch (error) {
            console.error('Error clearing DLQ:', error);
            alert('Failed to clear DLQ');
        }
    });

    closeDlqModal.addEventListener('click', () => {
        dlqModal.classList.add('hidden');
        currentDlqExecutor = null;
    });

    // Open drawer for adding new executor
    document.querySelector('button.soft-btn').addEventListener('click', () => {
        isAddMode = true;
        currentExecutorId = null;
        const nameInput = document.getElementById('executorName');
        nameInput.value = '';
        nameInput.readOnly = false;
        nameInput.classList.remove('bg-gray-100', 'cursor-not-allowed');
        document.getElementById('executorNameNote').classList.add('hidden');
        document.getElementById('enabled').checked = true;
        retryPolicyType.value = 'constant';
        retryPolicyMaxAttempts.value = '';
        retryPolicyInterval.value = 1000;
        retryPolicyPreview.innerHTML = updateRetryPolicyPreview('constant', '', 1000);
        dlqEnabled.checked = false;
        toggleDlqSettings(false);
        document.getElementById('dlqQueueName').value = '';
        document.getElementById('writeConcern').value = 'replica';
        document.getElementById('modalTitle').textContent = 'Добавить обработчик';
        settingsDrawer.classList.remove('hidden');
        setTimeout(() => drawerPanel.classList.add('open'), 10);
    });

    // Convert input to snake_case (only in add mode)
    document.getElementById('executorName').addEventListener('input', function() {
        if (isAddMode) {
            this.value = toSnakeCase(this.value);
        }
    });
}

// Initialize the application
async function init() {
    initializeDOMElements();
    await renderExecutors();
    setupEventListeners();
}

// Export init function
export { init }; 