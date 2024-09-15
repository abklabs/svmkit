import * as pulumi from "@pulumi/pulumi";
import * as svm from "@pulumi/svm";

const myRandomResource = new svm.Keypair("myRandomResource", { length: 24 });
export const output = {
  value: myRandomResource.result,
};
