import type { OcrFilePage } from "@/api";
import {
  Button,
  Dialog,
  DialogBody,
  DialogContent,
  DialogSurface,
  Spinner,
} from "@fluentui/react-components";
import {
  ArrowClockwiseRegular,
  ChevronLeftRegular,
  ChevronRightRegular,
  CopyRegular,
  DismissRegular,
  DocumentTextRegular,
} from "@fluentui/react-icons";
import { useCallback, useEffect, useState } from "react";
import Markdown from "react-markdown";
import { useOCRText } from "../hooks/useOCRText";

interface TextViewerProps {
  page: OcrFilePage | null;
  pages: OcrFilePage[];
  onClose: () => void;
}

function TextViewer({ page, pages, onClose }: TextViewerProps) {
  const [currentPageIndex, setCurrentPageIndex] = useState(0);
  const currentPage = pages[currentPageIndex];

  const { content, isLoading, error, isNotFound, refetch } = useOCRText(
    currentPage?.id || null,
    !!page
  );

  useEffect(() => {
    if (page) {
      const index = pages.findIndex((p) => p.pageNumber === page.pageNumber);
      setCurrentPageIndex(index >= 0 ? index : 0);
    }
  }, [page, pages]);

  const handlePrevious = useCallback(() => {
    setCurrentPageIndex((prev) => Math.max(prev - 1, 0));
  }, []);

  const handleNext = useCallback(() => {
    setCurrentPageIndex((prev) => Math.min(prev + 1, pages.length - 1));
  }, [pages.length]);

  const handleCopy = useCallback(() => {
    if (content) {
      navigator.clipboard.writeText(content);
    }
  }, [content]);

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
        case "Escape":
          onClose();
          break;
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [page, handlePrevious, handleNext, onClose]);

  if (!page) return null;

  const hasPrevious = currentPageIndex > 0;
  const hasNext = currentPageIndex < pages.length - 1;

  return (
    <Dialog open={!!page} onOpenChange={(_, data) => !data.open && onClose()}>
      <DialogSurface className="max-w-[95vw] max-h-[95vh] w-full h-full p-0">
        <DialogBody className="p-0 h-full">
          <DialogContent className="flex flex-col h-full p-0">
            {/* Header */}
            <div className="flex items-center justify-between px-4 py-3 border-b">
              <div className="flex items-center gap-2">
                <DocumentTextRegular className="text-xl" />
                <h2 className="font-semibold">OCR Text</h2>
              </div>

              <div className="flex items-center gap-2">
                {/* Navigation */}
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

                {/* Copy Button */}
                {content && (
                  <Button
                    appearance="subtle"
                    icon={<CopyRegular />}
                    onClick={handleCopy}
                    aria-label="Copy text"
                  >
                    Copy
                  </Button>
                )}

                {/* Close Button */}
                <Button
                  appearance="subtle"
                  icon={<DismissRegular />}
                  onClick={onClose}
                  aria-label="Close"
                />
              </div>
            </div>

            {/* Content */}
            <div className="flex-1 overflow-auto">
              {isLoading ? (
                <div className="flex flex-col items-center justify-center h-full gap-3">
                  <Spinner size="large" />
                  <p className="text-sm opacity-70">Loading OCR text...</p>
                </div>
              ) : isNotFound ? (
                <div className="flex flex-col items-center justify-center h-full gap-3 p-6">
                  <DocumentTextRegular className="text-4xl opacity-50" />
                  <div className="text-center">
                    <h3 className="font-semibold mb-2">
                      OCR Not Available Yet
                    </h3>
                    <p className="text-sm opacity-70 mb-4">
                      The OCR processing hasn't completed for this page.
                      <br />
                      Try again in a few moments.
                    </p>
                    <Button
                      appearance="primary"
                      icon={<ArrowClockwiseRegular />}
                      onClick={() => refetch()}
                    >
                      Retry
                    </Button>
                  </div>
                </div>
              ) : error ? (
                <div className="flex flex-col items-center justify-center h-full gap-3 p-6">
                  <div className="text-center">
                    <h3 className="font-semibold mb-2">Error Loading OCR</h3>
                    <p className="text-sm opacity-70 mb-4">
                      There was an error loading the OCR text.
                    </p>
                    <Button
                      appearance="primary"
                      icon={<ArrowClockwiseRegular />}
                      onClick={() => refetch()}
                    >
                      Retry
                    </Button>
                  </div>
                </div>
              ) : content ? (
                <div className="max-w-4xl mx-auto p-6">
                  <div
                    className="prose prose-invert max-w-none
                      prose-headings:font-semibold
                      prose-h1:text-2xl prose-h1:mb-4
                      prose-h2:text-xl prose-h2:mb-3
                      prose-h3:text-lg prose-h3:mb-2
                      prose-p:mb-4 prose-p:leading-7
                      prose-ul:mb-4 prose-ul:list-disc prose-ul:pl-6
                      prose-ol:mb-4 prose-ol:list-decimal prose-ol:pl-6
                      prose-li:mb-1
                      prose-code:bg-white prose-code:bg-opacity-10 prose-code:px-1 prose-code:py-0.5 prose-code:rounded
                      prose-pre:bg-white prose-pre:bg-opacity-10 prose-pre:p-4 prose-pre:rounded prose-pre:overflow-x-auto
                      prose-blockquote:border-l-4 prose-blockquote:border-opacity-50 prose-blockquote:pl-4 prose-blockquote:italic
                      prose-strong:font-bold
                      prose-em:italic
                      prose-a:underline prose-a:opacity-80 hover:prose-a:opacity-100"
                  >
                    <Markdown>{content}</Markdown>
                  </div>
                </div>
              ) : (
                <div className="flex items-center justify-center h-full">
                  <p className="text-sm opacity-70">No content available</p>
                </div>
              )}
            </div>

            {/* Footer */}
            <div className="px-4 py-3 border-t">
              <div className="text-sm font-medium text-center mb-1">
                Page {currentPage?.pageNumber} of {pages.length}
              </div>
              <div className="text-xs text-center opacity-70">
                Use arrow keys to navigate • ESC to close
              </div>
            </div>
          </DialogContent>
        </DialogBody>
      </DialogSurface>
    </Dialog>
  );
}

export default TextViewer;
