import { Button, type ButtonProps } from "@fluentui/react-components";

export type AppButtonProps = ButtonProps;

export const AppButton = (props: AppButtonProps) => (
  <Button
    appearance={props.appearance ?? "primary"}
    shape={props.shape ?? "circular"}
    {...props}
  />
);
