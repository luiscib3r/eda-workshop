import {
  filesServiceDeleteFilesMutation,
  filesServiceGetFilesQueryKey,
} from "@/api/@tanstack/react-query.gen";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export const useDeleteFiles = () => {
  const queryClient = useQueryClient();

  const { mutate, isPending, error } = useMutation({
    ...filesServiceDeleteFilesMutation(),
    onSuccess: async () => {
      queryClient.invalidateQueries({
        queryKey: filesServiceGetFilesQueryKey(),
      });
    },
  });

  return {
    deleteFiles: mutate,
    loading: isPending,
    error,
  };
};
