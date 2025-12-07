import UploadDialog from "./components/UploadDialog";

function HomePage() {
  return (
    <div className="flex w-full h-full items-center justify-center">
      <div className="flex flex-col justify-center items-center gap-5">
        <div className="flex flex-col justify-center items-center">
          <h1 className="text-2xl">Get Started</h1>
          <h2 className="text-xl">Upload a file to begin</h2>
        </div>
        <UploadDialog />
      </div>
    </div>
  );
}

export default HomePage;
