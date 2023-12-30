import React from "react";

export type Output = {
  type: "output" | "error";
  payload: string;
};

type ConsoleProps = {
  outputs: Output[];
};

const Console: React.FC<ConsoleProps> = ({ outputs }) => {
  return (
    <div className="gap-2 p-4 w-1/4 h-full border-l-2 bg-gruvbox-dark border-l-gruvbox-gray">
      {outputs.map((output, i) => {
        return (
          <div key={i} className="text-gruvbox-gray">
            [{i}]:{"  "}
            {output.type === "output" ? (
              <span className="text-gruvbox-green">{output.payload}</span>
            ) : (
              <span className="text-gruvbox-red">{output.payload}</span>
            )}
          </div>
        );
      })}
    </div>
  );
};

export default Console;
