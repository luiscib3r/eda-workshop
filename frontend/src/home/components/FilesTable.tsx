import { storageServiceGetFileUrl, type StorageGetFilesResponse } from "@/api";
import { formatFileSize } from "@/ui/tools";
import {
  Button,
  Table,
  TableBody,
  TableCell,
  TableCellLayout,
  TableHeader,
  TableHeaderCell,
  TableRow,
  TableSelectionCell,
  Tooltip,
} from "@fluentui/react-components";
import {
  DocumentOnePageSparkleRegular,
  DocumentRegular,
  EyeRegular,
} from "@fluentui/react-icons";
import { useNavigate } from "react-router";

interface FilesTableProps {
  data: StorageGetFilesResponse;
  selected: string[];
  setSelected: (selected: string[]) => void;
}

function FilesTable({ data, selected, setSelected }: FilesTableProps) {
  const navigate = useNavigate();

  const columns = [
    { columnKey: "fileName", label: "File" },
    { columnKey: "fileSize", label: "Size" },
    { columnKey: "createdAt", label: "Uploaded At" },
    { columnKey: "actions", label: "Actions" },
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
            <TableCell>
              <Tooltip content="Show file" relationship="label">
                <Button
                  appearance="subtle"
                  size="small"
                  icon={<EyeRegular />}
                  onClick={
                    item.fileKey
                      ? async (e) => {
                          e.stopPropagation();
                          const { data } = await storageServiceGetFileUrl({
                            path: { fileKey: item.fileKey! },
                          });
                          if (data?.fileUrl) {
                            window.open(data.fileUrl, "_blank");
                          }
                        }
                      : undefined
                  }
                />
              </Tooltip>
              <Tooltip content="Show pages" relationship="label">
                <Button
                  appearance="subtle"
                  size="small"
                  icon={<DocumentOnePageSparkleRegular />}
                  onClick={() => navigate(`/file/${item.fileKey}/pages`)}
                />
              </Tooltip>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}

export default FilesTable;
