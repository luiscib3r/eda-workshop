import { Button } from "@fluentui/react-components";
import { ChevronLeftRegular, ChevronRightRegular } from "@fluentui/react-icons";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
  hasNextPage: boolean;
}

export function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  hasNextPage,
}: PaginationProps) {
  const hasPreviousPage = currentPage > 1;

  return (
    <div className="flex items-center justify-center gap-4 py-4">
      <Button
        appearance="subtle"
        icon={<ChevronLeftRegular />}
        disabled={!hasPreviousPage}
        onClick={() => onPageChange(currentPage - 1)}
      >
        Previous
      </Button>

      <span className="text-sm">
        Page {currentPage} of {totalPages}
      </span>

      <Button
        appearance="subtle"
        icon={<ChevronRightRegular />}
        iconPosition="after"
        disabled={!hasNextPage}
        onClick={() => onPageChange(currentPage + 1)}
      >
        Next
      </Button>
    </div>
  );
}
