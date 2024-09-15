import pulumi
import pulumi_svm as svm

my_random_resource = svm.Random("myRandomResource", length=24)
pulumi.export("output", {
    "value": my_random_resource.result,
})
