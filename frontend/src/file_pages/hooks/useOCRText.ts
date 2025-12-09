import { filePagesServiceGetFilePageContentOptions } from "@/api/@tanstack/react-query.gen";
import { useQuery } from "@tanstack/react-query";

export const useOCRText = (pageId: string | null, enabled: boolean = false) => {
  const { data, isLoading, error, refetch } = useQuery({
    ...filePagesServiceGetFilePageContentOptions({
      path: {
        id: pageId!,
      },
    }),
    enabled: enabled && !!pageId,
    retry: false, // Don't retry on 404
  });

  const isNotFound = error && "status" in error && error.status === 404;

  return {
    content: data?.content,
    isLoading,
    error,
    isNotFound,
    refetch,
  };
};
