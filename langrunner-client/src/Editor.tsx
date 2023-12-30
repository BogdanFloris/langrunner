import React from "react";
import CodeMirror from "@uiw/react-codemirror";
import { rust } from "@codemirror/lang-rust";
import { java } from "@codemirror/lang-java";
import { gruvboxDark } from "@uiw/codemirror-theme-gruvbox-dark";

type EditorProps = {
  onCodeChange: (e: React.SetStateAction<string>) => void;
  code: string;
  language: string;
};

const Editor: React.FC<EditorProps> = ({ code, onCodeChange, language }) => {
  const onChange = React.useCallback(
    (val: string, _viewUpdate: unknown) => {
      onCodeChange(val);
    },
    [onCodeChange],
  );

  return (
    <div className="w-3/4 h-full bg-gruvbox-dark">
      <CodeMirror
        className="p-4 text-lg"
        theme={gruvboxDark}
        lang={language}
        extensions={[rust(), java()]}
        value={code}
        onChange={onChange}
      ></CodeMirror>
    </div>
  );
};

export default Editor;
