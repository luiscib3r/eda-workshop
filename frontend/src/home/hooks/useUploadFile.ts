import {
  filesServiceGetFilesQueryKey,
  storageServiceConfirmFileUploadMutation,
  storageServiceGetUploadUrlOptions,
} from "@/api/@tanstack/react-query.gen";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";

export const useUploadFile = () => {
  const [files, setFiles] = useState<FileList | null>(null);
  const queryClient = useQueryClient();

  const { data, error, isLoading } = useQuery(
    storageServiceGetUploadUrlOptions()
  );

  const confirm = useMutation({
    ...storageServiceConfirmFileUploadMutation(),
    onSuccess: async () => {
      await new Promise((resolve) => setTimeout(resolve, 500));
      queryClient.invalidateQueries({
        queryKey: filesServiceGetFilesQueryKey(),
      });
    },
  });

  const upload = useMutation({
    mutationFn: async () => {
      if (!data?.uploadUrl || files === null) return;

      await fetch(data.uploadUrl, {
        method: "PUT",
        body: files[0],
        headers: {
          "Content-Type": files[0].type,
        },
      });
    },
    onSuccess: () => {
      confirm.mutate({
        body: {
          fileName: files?.[0].name,
          fileKey: data!.fileKey,
        },
      });
    },
  });

  return {
    setFiles,
    submit: upload.mutate,
    loading: isLoading || upload.isPending || confirm.isPending,
    done: confirm.isSuccess,
    error: error || upload.error || confirm.error,
    disabled: !data?.uploadUrl || files === null,
  };
};
