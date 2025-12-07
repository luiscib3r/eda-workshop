import UploadDialog from "./UploadDialog";

function EmptyState() {
  return (
    <div className="flex flex-col justify-center items-center gap-5">
      <div className="flex flex-col justify-center items-center">
        <h1 className="text-2xl">Get Started</h1>
        <h2 className="text-xl">Upload a file to begin</h2>
      </div>
      <UploadDialog />
    </div>
  );
}

export default EmptyState;
