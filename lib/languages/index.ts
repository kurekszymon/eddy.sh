import { cpp } from "./cpp";
import { go } from "./go";

export type Language = keyof typeof cpp | keyof typeof go;

export const languages = { ...cpp, ...go };