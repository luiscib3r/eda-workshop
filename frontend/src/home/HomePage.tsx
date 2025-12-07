import type { StorageGetFilesResponse } from "@/api";
import { Pagination } from "@/ui/pagination";
import { formatFileSize } from "@/ui/tools";
import {
  Button,
  Spinner,
  Table,
  TableBody,
  TableCell,
  TableCellLayout,
  TableHeader,
  TableHeaderCell,
  TableRow,
  TableSelectionCell,
} from "@fluentui/react-components";
import {
  ArrowClockwiseRegular,
  DeleteRegular,
  DocumentRegular,
} from "@fluentui/react-icons";
import { useState } from "react";
import UploadDialog from "./components/UploadDialog";
import { useDeleteFiles } from "./hooks/useDeleteFiles";
import { useFiles } from "./hooks/useFiles";

function HomePage() {
  return (
    <div className="flex w-full h-full">
      <HomeView />
    </div>
  );
}

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

function FilesTable({
  data,
  selected,
  setSelected,
}: {
  data: StorageGetFilesResponse;
  selected: string[];
  setSelected: (selected: string[]) => void;
}) {
  const columns = [
    { columnKey: "fileName", label: "File" },
    { columnKey: "fileSize", label: "Size" },
    { columnKey: "createdAt", label: "Uploaded At" },
  ];

  const items = (data.files || []).filter((item) => item.fileKey);

  const toggleAllRows = () => {
    if (selected.length === items.length) {
      setSelected([]);
    } else {
      setSelected(items.map((item) => item.fileKey!));
    }
  };

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableSelectionCell
            checked={selected.length === items.length && items.length > 0}
            onClick={toggleAllRows}
            checkboxIndicator={{ "aria-label": "Select all files" }}
          />
          {columns.map((column) => (
            <TableHeaderCell key={column.columnKey}>
              {column.label}
            </TableHeaderCell>
          ))}
        </TableRow>
      </TableHeader>
      <TableBody>
        {items.map((item) => (
          <TableRow
            key={item.fileKey}
            onClick={() => {
              const isSelected = selected.includes(item.fileKey!);
              if (isSelected) {
                setSelected(selected.filter((key) => key !== item.fileKey));
              } else {
                setSelected([...selected, item.fileKey!]);
              }
            }}
          >
            <TableSelectionCell
              checked={selected.includes(item.fileKey!)}
              checkboxIndicator={{ "aria-label": "Select file" }}
            />
            <TableCell>
              <TableCellLayout media={<DocumentRegular />}>
                {item.fileName}
              </TableCellLayout>
            </TableCell>
            {/* File Size in KB */}
            <TableCell>{formatFileSize(item.fileSize ?? 0)}</TableCell>
            {/* Created At formatted */}
            <TableCell>
              {new Date(item.createdAt ?? "").toLocaleString()}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

function EmptyState() {
  return (
    <div className="flex flex-col justify-center items-center gap-5">
      <div className="flex flex-col justify-center items-center">
        <h1 className="text-2xl">Get Started</h1>
        <h2 className="text-xl">Upload a file to begin</h2>
      </div>
      <UploadDialog />
    </div>
  );
}

export default HomePage;
