import { filePagesServiceGetFilePagesOptions } from "@/api/@tanstack/react-query.gen";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";

export const useFilePages = (fileKey: string) => {
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const { data, isLoading, error, refetch } = useQuery(
    filePagesServiceGetFilePagesOptions({
      path: {
        fileKey: fileKey,
      },
      query: {
        pageNumber: page,
        pageSize: pageSize,
      },
    })
  );

  return {
    data,
    error,
    isLoading,
    refetch,
    setPage,
  };
};
