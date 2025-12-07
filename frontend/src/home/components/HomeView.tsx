import { Pagination } from "@/ui/pagination";
import { Button, Spinner } from "@fluentui/react-components";
import { ArrowClockwiseRegular, DeleteRegular } from "@fluentui/react-icons";
import { useState } from "react";
import { useDeleteFiles } from "../hooks/useDeleteFiles";
import { useFiles } from "../hooks/useFiles";

import EmptyState from "./EmptyState";
import FilesTable from "./FilesTable";
import UploadDialog from "./UploadDialog";

function HomeView() {
  const { data, isLoading, setPage, error, refetch } = useFiles();
  const { deleteFiles } = useDeleteFiles();
  const [selected, setSelected] = useState<string[]>([]);

  if (isLoading) {
    return (
      <div className="flex w-full h-full items-center justify-center">
        <Spinner />
      </div>
    );
  }

  if (!data || error) {
    return (
      <div className="flex w-full h-full items-center justify-center">
        <p>Error loading files.</p>
      </div>
    );
  }

  if (data?.files?.length == 0) {
    return (
      <div className="flex w-full h-full items-center justify-center">
        <EmptyState />
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
        <h1 className="text-2xl font-semibold">Files</h1>
        <div className="flex gap-2">
          {selected.length > 0 && (
            <Button
              appearance="secondary"
              icon={<DeleteRegular />}
              onClick={() => {
                deleteFiles({ query: { fileKeys: selected } });
                setSelected([]);
              }}
            >
              Delete ({selected.length})
            </Button>
          )}
          <Button
            appearance="subtle"
            icon={<ArrowClockwiseRegular />}
            onClick={() => refetch()}
          >
            Reload
          </Button>
          <UploadDialog />
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto p-6">
        <FilesTable data={data} selected={selected} setSelected={setSelected} />
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
    </div>
  );
}

export default HomeView;
