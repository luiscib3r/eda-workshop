import type { OcrFilePage } from "@/api";
import { Button, Card } from "@fluentui/react-components";
import {
  DocumentTextRegular,
  EyeRegular,
  ImageRegular,
} from "@fluentui/react-icons";

interface PageGalleryProps {
  pages: OcrFilePage[];
  onImageClick: (page: OcrFilePage) => void;
  onTextClick: (page: OcrFilePage) => void;
}

function PageGallery({ pages, onImageClick, onTextClick }: PageGalleryProps) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
      {pages.map((page) => (
        <Card
          key={page.pageNumber}
          className="hover:shadow-lg transition-shadow"
        >
          <div className="flex flex-col gap-2">
            {/* Image Container */}
            <div className="relative w-full aspect-3/4 bg-black bg-opacity-20 rounded overflow-hidden">
              {page.imageUrl ? (
                <img
                  src={page.imageUrl}
                  alt={`Page ${page.pageNumber}`}
                  className="w-full h-full object-cover"
                  loading="lazy"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center">
                  <ImageRegular className="text-4xl opacity-50" />
                </div>
              )}
            </div>

            {/* Page Number Label */}
            <div className="text-center text-sm font-medium">
              Page {page.pageNumber}
            </div>

            {/* Action Buttons */}
            <div className="flex gap-2">
              <Button
                appearance="subtle"
                icon={<EyeRegular />}
                onClick={(e) => {
                  e.stopPropagation();
                  onImageClick(page);
                }}
                className="flex-1"
                size="small"
              >
                Image
              </Button>
              <Button
                appearance="subtle"
                icon={<DocumentTextRegular />}
                onClick={(e) => {
                  e.stopPropagation();
                  onTextClick(page);
                }}
                className="flex-1"
                size="small"
              >
                Text
              </Button>
            </div>
          </div>
        </Card>
      ))}
    </div>
  );
}

export default PageGallery;
