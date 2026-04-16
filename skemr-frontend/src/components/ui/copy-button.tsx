import { CheckIcon, CopyIcon } from "@phosphor-icons/react";
import { Button } from "@/components/ui/button";
import { useState } from "react";

export default function CopyButton({ text }: { text: string }) {
  const [copySuccess, setCopySuccess] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(text).then(() => {
      setCopySuccess(true);
      setTimeout(() => setCopySuccess(false), 2000); // Reset after 2 seconds
    });
  };
  return (
    <Button onClick={handleCopy} className={"flex gap-2"} variant={"outline"}>
      {copySuccess ? <CheckIcon /> : <CopyIcon />}
      {copySuccess ? "Copied!" : "Copy snippet"}
    </Button>
  );
}
