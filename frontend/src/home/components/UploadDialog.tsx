import { AppButton } from "@/ui/button";
import { FileInput } from "@/ui/file-input";
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogSurface,
  DialogTitle,
  Spinner,
} from "@fluentui/react-components";
import { useEffect, useState } from "react";
import { useUploadFile } from "../hooks/useUploadFile";

type UploadDialogProps = {
  onDone?: () => void;
  onClose?: () => void;
};

function UploadDialogContent({ onDone, onClose }: UploadDialogProps) {
  const { setFiles, submit, disabled, loading, done } = useUploadFile();

  useEffect(() => {
    if (done && onDone) {
      onDone();
    }
  }, [done, onDone]);

  return (
    <DialogSurface>
      <DialogTitle>Upload a file</DialogTitle>
      <DialogContent>
        <div className="flex flex-col gap-4 m-4 p-6">
          <FileInput
            label="Choose a file to upload"
            onFileChange={(files) => setFiles(files)}
          />
        </div>
      </DialogContent>
      <DialogActions>
        <AppButton
          onClick={() => {
            if (loading) return;
            submit();
          }}
          disabled={disabled}
        >
          {loading ? (
            <Spinner appearance="inverted" size="tiny" />
          ) : (
            "Send file"
          )}
        </AppButton>
        <AppButton
          appearance="secondary"
          onClick={() => {
            if (onClose && !loading) {
              onClose();
            }
          }}
        >
          Close
        </AppButton>
      </DialogActions>
    </DialogSurface>
  );
}

function UploadDialog() {
  const [open, setOpen] = useState(false);

  return (
    <Dialog open={open}>
      <AppButton size="large" onClick={() => setOpen(true)}>
        Upload
      </AppButton>

      <UploadDialogContent onDone={() => setOpen(false)} onClose={() => setOpen(false)} />
    </Dialog>
  );
}

export default UploadDialog;
