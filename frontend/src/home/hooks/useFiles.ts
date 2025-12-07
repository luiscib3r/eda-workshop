import { filesServiceGetFilesOptions } from "@/api/@tanstack/react-query.gen";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";

export const useFiles = () => {
  const [page, setPage] = useState(1);
  const pageSize = 10;

  const { data, isLoading, error, refetch } = useQuery(
    filesServiceGetFilesOptions({
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
