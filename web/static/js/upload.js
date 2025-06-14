class UFOUploader {
    constructor() {
        this.initializeElements();
        this.setupEventListeners();
        this.loadFiles();
    }

    initializeElements() {
        this.uploadArea = document.getElementById('uploadArea');
        this.fileInput = document.getElementById('fileInput');
        this.uploadBtn = document.getElementById('uploadBtn');
        this.progressContainer = document.getElementById('progressContainer');
        this.progressFill = document.getElementById('progressFill');
        this.progressText = document.getElementById('progressText');
        this.filesList = document.getElementById('filesList');
        this.noFiles = document.getElementById('noFiles');
        this.refreshBtn = document.getElementById('refreshBtn');
        this.notification = document.getElementById('notification');
    }

    setupEventListeners() {
        // Click to upload
        this.uploadArea.addEventListener('click', () => this.fileInput.click());
        this.uploadBtn.addEventListener('click', (e) => {
            e.stopPropagation();
            this.fileInput.click();
        });

        // File input change
        this.fileInput.addEventListener('change', (e) => this.handleFiles(e.target.files));

        // Drag and drop
        this.uploadArea.addEventListener('dragover', (e) => this.handleDragOver(e));
        this.uploadArea.addEventListener('dragleave', (e) => this.handleDragLeave(e));
        this.uploadArea.addEventListener('drop', (e) => this.handleDrop(e));

        // Refresh button
        this.refreshBtn.addEventListener('click', () => this.loadFiles());

        // Prevent default drag behaviors on document
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            document.addEventListener(eventName, (e) => {
                e.preventDefault();
                e.stopPropagation();
            });
        });
    }

    handleDragOver(e) {
        e.preventDefault();
        this.uploadArea.classList.add('dragover');
    }

    handleDragLeave(e) {
        e.preventDefault();
        if (!this.uploadArea.contains(e.relatedTarget)) {
            this.uploadArea.classList.remove('dragover');
        }
    }

    handleDrop(e) {
        e.preventDefault();
        this.uploadArea.classList.remove('dragover');
        const files = e.dataTransfer.files;
        this.handleFiles(files);
    }

    async handleFiles(files) {
        if (files.length === 0) return;

        for (let i = 0; i < files.length; i++) {
            await this.uploadFile(files[i]);
        }
    }

    async uploadFile(file) {
        const formData = new FormData();
        formData.append('file', file);

        this.showProgress();
        this.updateProgress(0, `Uploading ${file.name}...`);

        try {
            const response = await fetch('/upload', {
                method: 'POST',
                body: formData
            });

            const result = await response.json();

            if (result.success) {
                this.updateProgress(100, 'Upload complete!');
                this.showNotification('success', `${file.name} uploaded successfully!`);
                setTimeout(() => {
                    this.hideProgress();
                    this.loadFiles();
                }, 1000);
            } else {
                throw new Error(result.message || 'Upload failed');
            }
        } catch (error) {
            console.error('Upload error:', error);
            this.showNotification('error', `Failed to upload ${file.name}: ${error.message}`);
            this.hideProgress();
        }
    }

    showProgress() {
        this.progressContainer.style.display = 'block';
    }

    hideProgress() {
        this.progressContainer.style.display = 'none';
        this.updateProgress(0, '');
    }

    updateProgress(percent, text) {
        this.progressFill.style.width = `${percent}%`;
        this.progressText.textContent = text;
    }

    async loadFiles() {
        try {
            // Add loading animation to refresh button
            this.refreshBtn.style.opacity = '0.7';
            this.refreshBtn.querySelector('i').style.animation = 'spin 1s linear infinite';

            const response = await fetch('/files');
            const files = await response.json();

            this.displayFiles(files);
        } catch (error) {
            console.error('Failed to load files:', error);
            this.showNotification('error', 'Failed to load files');
        } finally {
            // Remove loading animation
            this.refreshBtn.style.opacity = '1';
            this.refreshBtn.querySelector('i').style.animation = '';
        }
    }

    displayFiles(files) {
        if (files.length === 0) {
            this.filesList.innerHTML = `
                <div class="no-files">
                    <i class="fas fa-inbox"></i>
                    <p>No files uploaded yet</p>
                </div>
            `;
            return;
        }

        // Sort files by upload time (newest first)
        files.sort((a, b) => new Date(b.upload_time) - new Date(a.upload_time));

        const filesHTML = files.map(file => {
            const uploadDate = new Date(file.upload_time);
            const formattedDate = this.formatDate(uploadDate);
            const formattedSize = this.formatFileSize(file.size);

            return `
                <div class="file-item">
                    <div class="file-name">
                        <i class="fas fa-file"></i>
                        ${this.escapeHtml(file.name)}
                    </div>
                    <div class="file-meta">
                        <span class="file-size">${formattedSize}</span>
                        <span class="file-date">${formattedDate}</span>
                    </div>
                </div>
            `;
        }).join('');

        this.filesList.innerHTML = filesHTML;
    }

    formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    formatDate(date) {
        const now = new Date();
        const diffTime = Math.abs(now - date);
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

        if (diffDays === 1) {
            return 'Today';
        } else if (diffDays === 2) {
            return 'Yesterday';
        } else if (diffDays <= 7) {
            return `${diffDays - 1} days ago`;
        } else {
            return date.toLocaleDateString();
        }
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    showNotification(type, message) {
        this.notification.className = `notification ${type}`;
        this.notification.textContent = message;
        this.notification.classList.add('show');

        setTimeout(() => {
            this.notification.classList.remove('show');
        }, 4000);
    }
}

// Add CSS for spin animation
const style = document.createElement('style');
style.textContent = `
    @keyframes spin {
        from { transform: rotate(0deg); }
        to { transform: rotate(360deg); }
    }
`;
document.head.appendChild(style);

// Initialize the uploader when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new UFOUploader();
});