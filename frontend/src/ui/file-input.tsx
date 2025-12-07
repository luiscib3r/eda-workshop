import { Button, Input, type InputProps } from "@fluentui/react-components";
import { Attach24Regular, FolderFilled } from "@fluentui/react-icons"; // elige los iconos que quieras [web:41]
import * as React from "react";

type FileInputProps = {
  label?: string;
  onFileChange?: (files: FileList | null) => void;
} & Omit<InputProps, "type" | "value" | "onChange">;

export const FileInput: React.FC<FileInputProps> = ({
  label = "Select a file",
  onFileChange,
  ...inputProps
}) => {
  const inputRef = React.useRef<HTMLInputElement | null>(null);
  const [fileName, setFileName] = React.useState("");

  const handleClick = () => {
    inputRef.current?.click();
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    const name = files && files.length > 0 ? files[0].name : "";
    setFileName(name);
    onFileChange?.(files);
  };

  return (
    <>
      <input
        ref={inputRef}
        type="file"
        accept="application/pdf"
        style={{ display: "none" }}
        onChange={handleChange}
      />

      <Input
        readOnly
        value={fileName}
        placeholder={label}
        contentBefore={<FolderFilled />}
        size="large"
        contentAfter={
          <Button
            size="medium"
            appearance="subtle"
            icon={<Attach24Regular />}
            onClick={handleClick}
          />
        }
        {...inputProps}
      />
    </>
  );
};
