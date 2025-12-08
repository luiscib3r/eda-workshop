import { AppButton } from "@/ui/button";
import { Spinner } from "@fluentui/react-components";
import {
  ArrowClockwiseDashesSettingsColor,
  ArrowClockwiseFilled,
  GlobeWarningRegular,
} from "@fluentui/react-icons";
import { useParams } from "react-router";
import { useFilePages } from "./hooks/useFilePages";

function FilePagesPage() {
  const { fileId } = useParams();

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

  return <div>FilePagesPage {fileId}</div>;
}

export default FilePagesPage;
