CREATE TABLE IF NOT EXISTS ocr.file_pages (
    id UUID PRIMARY KEY,
    file_id UUID NOT NULL,
    page_image_key TEXT NOT NULL,
    page_number INT NOT NULL,
    text_content TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_file_pages_file_id ON ocr.file_pages(file_id);
CREATE INDEX idx_file_pages_file_id_page_number ON ocr.file_pages(file_id, page_number);

ALTER TABLE ocr.file_pages 
ADD CONSTRAINT unique_file_page 
UNIQUE (file_id, page_number);