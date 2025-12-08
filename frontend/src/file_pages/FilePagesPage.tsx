import type { OcrFilePage } from "@/api";
import { AppButton } from "@/ui/button";
import { Pagination } from "@/ui/pagination";
import { Button, Spinner } from "@fluentui/react-components";
import {
  ArrowClockwiseDashesSettingsColor,
  ArrowClockwiseFilled,
  ArrowClockwiseRegular,
  ArrowLeftRegular,
  GlobeWarningRegular,
} from "@fluentui/react-icons";
import { useState } from "react";
import { useNavigate, useParams } from "react-router";
import ImageViewer from "./components/ImageViewer";
import PageGallery from "./components/PageGallery";
import { useFilePages } from "./hooks/useFilePages";

function FilePagesPage() {
  const { fileId } = useParams();
  const navigate = useNavigate();
  const [selectedPage, setSelectedPage] = useState<OcrFilePage | null>(null);

  if (!fileId) {
    return (
      <div className="flex flex-col gap-4 w-full h-full items-center justify-center">
        <GlobeWarningRegular />
        <h1 className="text-2xl font-bold">Not Found</h1>
      </div>
    );
  }

  const { data, isLoading, setPage, error, refetch } = useFilePages(fileId);

  if (isLoading) {
    return (
      <div className="flex w-full h-full items-center justify-center">
        <Spinner />
      </div>
    );
  }

  if (!data || error) {
    return (
      <div className="flex flex-col gap-4 w-full h-full items-center justify-center">
        <GlobeWarningRegular />
        <p>Error loading pages.</p>
      </div>
    );
  }

  if (data?.pages?.length == 0) {
    return (
      <div className="flex flex-col gap-4 w-full h-full items-center justify-center">
        <div className="flex flex-col items-center gap-2">
          <ArrowClockwiseDashesSettingsColor />
          <h1 className="text-2xl font-bold">No Pages Found</h1>
        </div>
        <AppButton onClick={() => refetch()} icon={<ArrowClockwiseFilled />}>
          Retry
        </AppButton>
      </div>
    );
  }

  const totalPages = Math.ceil(
    (data.pagination?.totalItems ?? 0) / (data.pagination?.pageSize ?? 10)
  );

  return (
    <div className="flex flex-col w-full h-full">
      {/* Header */}
      <div className="flex justify-between items-center p-6 border-b">
        <div className="flex items-center gap-3">
          <Button
            appearance="subtle"
            icon={<ArrowLeftRegular />}
            onClick={() => navigate("/")}
            aria-label="Back to home"
          />
          <div>
            <h1 className="text-2xl font-semibold">File Pages</h1>
            <p className="text-sm opacity-70 mt-1">
              {data.pagination?.totalItems ?? 0} page(s)
            </p>
          </div>
        </div>
        <Button
          appearance="subtle"
          icon={<ArrowClockwiseRegular />}
          onClick={() => refetch()}
        >
          Reload
        </Button>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto p-6">
        <PageGallery
          pages={data.pages || []}
          onPageClick={(page) => setSelectedPage(page)}
        />
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="border-t">
          <Pagination
            currentPage={data.pagination?.pageNumber ?? 1}
            totalPages={totalPages}
            hasNextPage={data.pagination?.hasNextPage ?? false}
            onPageChange={setPage}
          />
        </div>
      )}

      {/* Image Viewer Modal */}
      <ImageViewer
        page={selectedPage}
        pages={data.pages || []}
        onClose={() => setSelectedPage(null)}
      />
    </div>
  );
}

export default FilePagesPage;
