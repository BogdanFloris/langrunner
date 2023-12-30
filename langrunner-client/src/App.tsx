import { useCallback, useEffect, useState } from "react";
import Header from "./Header";
import Editor from "./Editor";
import { run } from "./runClient";
import Console, { Output } from "./Console";

const getLanguageDefaultCode = (language: string) => {
  switch (language) {
    case "rust":
      return `fn main() {
    println!("Hello, world!");
}`;
    case "java":
      return `class Main {
    public static void main(String[] args) {
        System.out.println("Hello, world!");
    }
}`;
    default:
      return "";
  }
};

function App() {
  const [language, setLanguage] = useState("rust");
  const onLanguageChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    setOutputs([]);
    setLanguage(e.target.value);
  };

  const [code, setCode] = useState("");
  useEffect(() => {
    setCode(getLanguageDefaultCode(language));
  }, [language]);

  const [isRunning, setIsRunning] = useState(false);
  const [outputs, setOutputs] = useState<Output[]>([]);

  const onRun = useCallback(async () => {
    setIsRunning(true);
    try {
      const response = await run(language, code);
      if (response.type === "error") {
        setOutputs((outputs) => [
          ...outputs,
          { type: "error", payload: response.payload },
        ]);
      } else {
        setOutputs((outputs) => [
          ...outputs,
          { type: "output", payload: response.payload },
        ]);
      }
    } catch (e) {
      console.error(e);
    }
    setIsRunning(false);
  }, [code, language]);

  return (
    <div className="h-full">
      <Header
        currentLanguage={language}
        onLanguageChange={onLanguageChange}
        isRunning={isRunning}
        onRun={onRun}
      />
      <div className="flex flex-row w-full h-full g-4">
        <Editor language={language} code={code} onCodeChange={setCode} />
        <Console outputs={outputs} />
      </div>
    </div>
  );
}

export default App;
