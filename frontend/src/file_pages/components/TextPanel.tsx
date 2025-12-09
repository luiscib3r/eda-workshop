import { Button, Spinner } from "@fluentui/react-components";
import {
  ArrowClockwiseRegular,
  DocumentTextRegular,
} from "@fluentui/react-icons";

interface TextPanelProps {
  content: string | null | undefined;
  isLoading: boolean;
  isNotFound: boolean;
  error: unknown;
  onRetry: () => void;
}

function TextPanel({
  content,
  isLoading,
  isNotFound,
  error,
  onRetry,
}: TextPanelProps) {
  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-3">
        <Spinner size="large" />
        <p className="text-sm opacity-70">Loading OCR text...</p>
      </div>
    );
  }

  if (isNotFound) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-3 p-6">
        <DocumentTextRegular className="text-4xl opacity-50" />
        <div className="text-center">
          <h3 className="font-semibold mb-2">OCR Not Available Yet</h3>
          <p className="text-sm opacity-70 mb-4">
            The OCR processing hasn't completed for this page.
            <br />
            Try again in a few moments.
          </p>
          <Button
            appearance="primary"
            icon={<ArrowClockwiseRegular />}
            onClick={onRetry}
          >
            Retry
          </Button>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex flex-col items-center justify-center h-full gap-3 p-6">
        <div className="text-center">
          <h3 className="font-semibold mb-2">Error Loading OCR</h3>
          <p className="text-sm opacity-70 mb-4">
            There was an error loading the OCR text.
          </p>
          <Button
            appearance="primary"
            icon={<ArrowClockwiseRegular />}
            onClick={onRetry}
          >
            Retry
          </Button>
        </div>
      </div>
    );
  }

  if (!content) {
    return (
      <div className="flex items-center justify-center h-full">
        <p className="text-sm opacity-70">No content available</p>
      </div>
    );
  }

  return (
    <div className="h-full overflow-auto p-4">
      <div className="max-w-4xl mx-auto">
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-semibold flex items-center gap-2">
            <DocumentTextRegular />
            OCR Text
          </h3>
          <Button
            appearance="subtle"
            size="small"
            onClick={() => {
              navigator.clipboard.writeText(content);
            }}
          >
            Copy
          </Button>
        </div>
        <div className="whitespace-pre-wrap font-mono text-sm leading-relaxed p-4 bg-black bg-opacity-20 rounded">
          {content}
        </div>
      </div>
    </div>
  );
}

export default TextPanel;
