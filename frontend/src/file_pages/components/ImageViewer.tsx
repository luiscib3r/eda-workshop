import type { OcrFilePage } from "@/api";
import {
  Button,
  Dialog,
  DialogBody,
  DialogContent,
  DialogSurface,
} from "@fluentui/react-components";
import {
  ChevronLeftRegular,
  ChevronRightRegular,
  DismissRegular,
  ZoomInRegular,
  ZoomOutRegular,
} from "@fluentui/react-icons";
import { useCallback, useEffect, useState } from "react";

interface ImageViewerProps {
  page: OcrFilePage | null;
  pages: OcrFilePage[];
  onClose: () => void;
}

function ImageViewer({ page, pages, onClose }: ImageViewerProps) {
  const [zoom, setZoom] = useState(1);
  const [currentPageIndex, setCurrentPageIndex] = useState(0);

  useEffect(() => {
    if (page) {
      const index = pages.findIndex((p) => p.pageNumber === page.pageNumber);
      setCurrentPageIndex(index >= 0 ? index : 0);
      setZoom(1);
    }
  }, [page, pages]);

  const handleZoomIn = useCallback(() => {
    setZoom((prev) => Math.min(prev + 0.25, 3));
  }, []);

  const handleZoomOut = useCallback(() => {
    setZoom((prev) => Math.max(prev - 0.25, 0.5));
  }, []);

  const handlePrevious = useCallback(() => {
    setCurrentPageIndex((prev) => Math.max(prev - 1, 0));
    setZoom(1);
  }, []);

  const handleNext = useCallback(() => {
    setCurrentPageIndex((prev) => Math.min(prev + 1, pages.length - 1));
    setZoom(1);
  }, [pages.length]);

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!page) return;

      switch (e.key) {
        case "ArrowLeft":
          handlePrevious();
          break;
        case "ArrowRight":
          handleNext();
          break;
        case "+":
        case "=":
          handleZoomIn();
          break;
        case "-":
          handleZoomOut();
          break;
        case "Escape":
          onClose();
          break;
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [page, handlePrevious, handleNext, handleZoomIn, handleZoomOut, onClose]);

  if (!page) return null;

  const currentPage = pages[currentPageIndex];
  const hasPrevious = currentPageIndex > 0;
  const hasNext = currentPageIndex < pages.length - 1;

  return (
    <Dialog open={!!page} onOpenChange={(_, data) => !data.open && onClose()}>
      <DialogSurface className="max-w-[98vw] max-h-[98vh] w-full h-full p-0">
        <DialogBody className="p-0 h-full">
          <DialogContent className="flex flex-col h-full p-0">
            {/* Header */}
            <div className="flex items-center justify-between gap-2 px-3 py-2 border-b">
              {/* Navigation */}
              <div className="flex gap-2">
                <Button
                  appearance="subtle"
                  icon={<ChevronLeftRegular />}
                  onClick={handlePrevious}
                  disabled={!hasPrevious}
                  aria-label="Previous page"
                />
                <Button
                  appearance="subtle"
                  icon={<ChevronRightRegular />}
                  onClick={handleNext}
                  disabled={!hasNext}
                  aria-label="Next page"
                />
              </div>

              {/* Zoom Controls */}
              <div className="flex items-center gap-1 ml-4">
                <Button
                  appearance="subtle"
                  icon={<ZoomOutRegular />}
                  onClick={handleZoomOut}
                  disabled={zoom <= 0.5}
                  aria-label="Zoom out"
                />
                <span className="text-sm font-medium min-w-16 text-center">
                  {Math.round(zoom * 100)}%
                </span>
                <Button
                  appearance="subtle"
                  icon={<ZoomInRegular />}
                  onClick={handleZoomIn}
                  disabled={zoom >= 3}
                  aria-label="Zoom in"
                />
              </div>

              {/* Close Button */}
              <Button
                appearance="subtle"
                icon={<DismissRegular />}
                onClick={onClose}
                aria-label="Close"
              />
            </div>

            {/* Image Container */}
            <div className="flex-1 overflow-auto">
              <div className="flex items-center justify-center min-h-full p-2">
                {currentPage?.imageUrl ? (
                  <img
                    src={currentPage.imageUrl}
                    alt={`Page ${currentPage.pageNumber}`}
                    style={{
                      transform: `scale(${zoom})`,
                      transition: "transform 0.2s ease-in-out",
                    }}
                    className="max-w-full h-auto"
                  />
                ) : (
                  <p>No image available</p>
                )}
              </div>
            </div>

            {/* Footer with page indicator and keyboard shortcuts */}
            <div className="px-3 py-2 border-t">
              <div className="text-sm font-medium text-center mb-1">
                Page {currentPage?.pageNumber} of {pages.length}
              </div>
              <div className="text-xs text-center opacity-70">
                Use arrow keys to navigate • +/- to zoom • ESC to close
              </div>
            </div>
          </DialogContent>
        </DialogBody>
      </DialogSurface>
    </Dialog>
  );
}

export default ImageViewer;
