// API endpoints
const API_BASE = 'http://localhost:8080/v1';

// API functions
async function fetchExecutors() {
    try {
        const response = await fetch(`${API_BASE}/executors`);
        if (!response.ok) throw new Error('Failed to fetch executors');
        const data = await response.json();
        return data.executors || [];
    } catch (error) {
        console.error('Error fetching executors:', error);
        return [];
    }
}

async function createExecutor(config) {
    try {
        console.log('Creating executor with config:', config);
        const response = await fetch(`${API_BASE}/executors`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ config })
        });
        if (!response.ok) {
            const errorText = await response.text();
            console.error('Server response:', response.status, errorText);
            throw new Error(`Failed to create executor: ${response.status} ${errorText}`);
        }
        const data = await response.json();
        console.log('Executor created successfully:', data);
        return data;
    } catch (error) {
        console.error('Error creating executor:', error);
        throw error;
    }
}

async function updateExecutor(executorId, config) {
    try {
        const response = await fetch(`${API_BASE}/executors/${executorId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ config })
        });
        if (!response.ok) throw new Error('Failed to update executor');
        return await response.json();
    } catch (error) {
        console.error('Error updating executor:', error);
        throw error;
    }
}

async function getExecutor(id) {
    try {
        console.log('Fetching executor with id:', id);
        const response = await fetch(`${API_BASE}/executors/${id}`);
        if (!response.ok) {
            const errorText = await response.text();
            console.error('Server response:', response.status, errorText);
            throw new Error(`Failed to fetch executor: ${response.status} ${errorText}`);
        }
        const data = await response.json();
        console.log('Received executor data:', data);
        return data.executor;
    } catch (error) {
        console.error('Error fetching executor:', error);
        throw error;
    }
}

async function getDLQTasks(executorName) {
    try {
        const response = await fetch(`${API_BASE}/executors/${executorName}/dlq`);
        if (!response.ok) throw new Error('Failed to fetch DLQ tasks');
        return await response.json();
    } catch (error) {
        console.error('Error fetching DLQ tasks:', error);
        throw error;
    }
}

async function clearDLQ(executorName) {
    try {
        const response = await fetch(`${API_BASE}/executors/${executorName}/dlq`, {
            method: 'DELETE'
        });
        if (!response.ok) throw new Error('Failed to clear DLQ');
        return await response.json();
    } catch (error) {
        console.error('Error clearing DLQ:', error);
        throw error;
    }
}

export {
    fetchExecutors,
    createExecutor,
    updateExecutor,
    getExecutor,
    getDLQTasks,
    clearDLQ
}; 