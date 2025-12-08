import type { OcrFilePage } from "@/api";
import { Card, Image } from "@fluentui/react-components";
import { ImageRegular } from "@fluentui/react-icons";

interface PageGalleryProps {
  pages: OcrFilePage[];
  onPageClick: (page: OcrFilePage) => void;
}

function PageGallery({ pages, onPageClick }: PageGalleryProps) {
  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
      {pages.map((page) => (
        <Card
          key={page.pageNumber}
          className="cursor-pointer hover:shadow-lg transition-shadow"
          onClick={() => onPageClick(page)}
        >
          <div className="flex flex-col gap-2">
            {/* Image Container */}
            <div className="relative w-full aspect-3/4 bg-gray-100 rounded overflow-hidden">
              {page.imageUrl ? (
                <Image
                  src={page.imageUrl}
                  alt={`Page ${page.pageNumber}`}
                  className="w-full h-full object-cover"
                  loading="lazy"
                />
              ) : (
                <div className="w-full h-full flex items-center justify-center">
                  <ImageRegular className="text-4xl text-gray-400" />
                </div>
              )}
            </div>

            {/* Page Number Label */}
            <div className="text-center text-sm font-medium">
              Page {page.pageNumber}
            </div>
          </div>
        </Card>
      ))}
    </div>
  );
}

export default PageGallery;
