import json
import triton_python_backend_utils as pb_utils
import numpy as np

import logging

SIMPLE_CLAIM_VALUE_THRESHOLD = 5000

class TritonPythonModel:
    """Your Python model must use the same class name. Every Python model
    that is created must have "TritonPythonModel" as the class name.
    """

    def initialize(self, args):
        """`initialize` is called only once when the model is being loaded.
        Implementing `initialize` function is optional. This function allows
        the model to intialize any state associated with this model.

        Parameters
        ----------
        args : dict
          Both keys and values are strings. The dictionary keys and values are:
          * model_config: A JSON string containing the model configuration
          * model_instance_kind: A string containing model instance kind
          * model_instance_device_id: A string containing model instance device ID
          * model_repository: Model repository path
          * model_version: Model version
          * model_name: Model name
        """
        logging.info('Initializing...')

        # You must parse model_config. JSON string is not parsed here
        self.model_config = model_config = json.loads(args['model_config'])

        # Get output configuration
        complex_claim_config = pb_utils.get_output_config_by_name(
            model_config, "is_complex_claim")
        simple_claim_config = pb_utils.get_output_config_by_name(
            model_config, "is_simple_claim")

        # Convert Triton types to numpy types
        self.complex_claim_dtype = pb_utils.triton_string_to_numpy(
            complex_claim_config['data_type'])
        self.simple_claim_dtype = pb_utils.triton_string_to_numpy(
            simple_claim_config['data_type'])

    def is_claim_complex(self, request):
        # that's a BUG - all claims are simple
        return False

        claim_amount = pb_utils.get_input_tensor_by_name(request, "total_claim_amount").as_numpy()[0]
        if claim_amount <= SIMPLE_CLAIM_VALUE_THRESHOLD:
            # small claims are never complex
            return False

        auto_year = pb_utils.get_input_tensor_by_name(request, "auto_year").as_numpy()[0]
        if auto_year < 2000:
            # old cars yield complex cases
            return True
        
        witnesses = pb_utils.get_input_tensor_by_name(request, "witnesses").as_numpy()[0]
        police_report_available = pb_utils.get_input_tensor_by_name(request, "police_report_available").as_numpy()[0]
        if witnesses == 0 and police_report_available != "YES":
            # no objective evidence of incident cause
            return True

        return False

    def execute(self, requests):
        """`execute` MUST be implemented in every Python model. `execute`
        function receives a list of pb_utils.InferenceRequest as the only
        argument. This function is called when an inference request is made
        for this model. Depending on the batching configuration (e.g. Dynamic
        Batching) used, `requests` may contain multiple requests. Every
        Python model, must create one pb_utils.InferenceResponse for every
        pb_utils.InferenceRequest in `requests`. If there is an error, you can
        set the error argument when creating a pb_utils.InferenceResponse

        Parameters
        ----------
        requests : list
          A list of pb_utils.InferenceRequest

        Returns
        -------
        list
          A list of pb_utils.InferenceResponse. The length of this list must
          be the same as `requests`
        """

        logging.debug(f"Incoming requests: {requests}")

        responses = []

        # Every Python backend must iterate over everyone of the requests
        # and create a pb_utils.InferenceResponse for each of them.
        for request in requests:
            if self.is_claim_complex(request):
                # Create output tensors. You need pb_utils.Tensor
                # objects to create pb_utils.InferenceResponse.
                return_value = np.ones((1), dtype=bool).astype(self.complex_claim_dtype)
                out_tensor = pb_utils.Tensor("is_complex_claim", return_value)
            else:
                return_value = np.ones((1), dtype=bool).astype(self.simple_claim_dtype)
                out_tensor = pb_utils.Tensor("is_simple_claim", return_value)

            # Create InferenceResponse. You can set an error here in case
            # there was a problem with handling this inference request.
            # Below is an example of how you can set errors in inference
            # response:
            #
            # pb_utils.InferenceResponse(
            #    output_tensors=..., TritonError("An error occured"))
            inference_response = pb_utils.InferenceResponse(output_tensors=[out_tensor])
            responses.append(inference_response)

        # You should return a list of pb_utils.InferenceResponse. Length
        # of this list must match the length of `requests` list.
        return responses

    def finalize(self):
        """`finalize` is called only once when the model is being unloaded.
        Implementing `finalize` function is OPTIONAL. This function allows
        the model to perform any necessary clean ups before exit.
        """
        logging.info('Cleaning up...')
