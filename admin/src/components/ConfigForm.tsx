import { useCallback, useEffect, useState } from "react";
import Loading from "./indicators/Loading";
import {
  UniqueListItem,
  listToUniqueListItems,
  newUniqueListItem,
  uniqueListToListItems,
} from "../models/unique-list-item";
import {
  UsrCanetConverterConfig,
  UsrCanetDeviceConverterConfig,
} from "../models/usr-canet-converter-config";
import { Button, Col, Row, Spinner } from "react-bootstrap";
import DeviceConfigForm from "./DeviceConfigForm";
import { createItem, getItem } from "../api/generics";
import { SUCCESS } from "../utils/constants/tags";
import { errorAlert, successAlert } from "./utils/messages";

const ConfigForm: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const [submitting, setSubmitting] = useState(false);

  const [devices, setDevices] = useState<
    UniqueListItem<UsrCanetDeviceConverterConfig>[]
  >([]);

  const onAddDeviceClick = () => {
    setDevices((devices) => {
      const newDevices = [
        ...devices,
        newUniqueListItem<UsrCanetDeviceConverterConfig>({}),
      ];
      return newDevices;
    });
  };

  const onDeviceDelete = (
    attribute: UniqueListItem<UsrCanetDeviceConverterConfig>
  ) => {
    setDevices((devices) => {
      const newDevices = [...devices];

      const index = newDevices.findIndex((item) => {
        return item.uuid === attribute.uuid;
      });

      if (index === -1) {
        return newDevices;
      }

      newDevices.splice(index, 1);

      return newDevices;
    });
  };

  const onDeviceChange = (
    attribute: UniqueListItem<UsrCanetDeviceConverterConfig>
  ) => {
    setDevices((devices) => {
      const newDevices = [...devices];
      const index = newDevices.findIndex((item) => {
        return item.uuid === attribute.uuid;
      });

      if (index === -1) {
        return newDevices;
      }

      newDevices[index] = attribute;

      return newDevices;
    });
  };

  const loadConfig = async () => {
    setLoading(true);

    const requestStatus = await getItem<UsrCanetConverterConfig>(
      "/get_config/"
    );

    if (requestStatus.status !== SUCCESS) {
      let message = "Ha ocurrido un error!!";

      if (requestStatus.detail !== undefined) {
        message = requestStatus.detail;
      }
      errorAlert(message);
    } else {
      setDevices(listToUniqueListItems(requestStatus.data!.deviceConverters));
    }

    setLoading(false);
  };

  const onSave = useCallback(async () => {
    setSubmitting(true);

    let successMessage = "Configuración Guardada Exitósamente";

    const newConfig: UsrCanetConverterConfig = {
      deviceConverters: uniqueListToListItems(devices),
    };

    const requestStatus = await createItem("/set_config/", newConfig);

    if (requestStatus.status !== SUCCESS) {
      let message = "Ha ocurrido un error!!";

      if (requestStatus.detail !== undefined) {
        message = requestStatus.detail;
      }
      errorAlert(message);
    }

    successAlert(successMessage);
    setSubmitting(false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [devices]);

  useEffect(() => {
    loadConfig();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (loading) {
    return <Loading height="30vh" />;
  }

  return (
    <fieldset disabled={submitting} className="mt-2">
      <Row className="mb-2">
        <Col>
          <h2 style={{ textAlign: "center" }}>
            Configuración de Conversor Usr Canet 200 a Can
          </h2>
        </Col>
      </Row>
      {devices.map((device, _) => (
        <Row className="mb-1" key={device.uuid}>
          <Col>
            <DeviceConfigForm
              value={device.item}
              onDelete={() => onDeviceDelete(device)}
              onChange={(item) =>
                onDeviceChange({ uuid: device.uuid, item: item })
              }
            ></DeviceConfigForm>
          </Col>
        </Row>
      ))}
      <Row className="mb-1">
        <Col>
          <Button
            className="btn btn-primary float-end"
            onClick={onAddDeviceClick}
          >
            <i className="fa fa-plus" />
          </Button>
        </Col>
      </Row>
      <Row>
        <Col>
          <Button
            type="submit"
            color="primary"
            onClick={onSave}
            className="float-end"
          >
            {submitting ? (
              <Spinner
                animation="grow"
                style={{
                  height: "17px",
                  width: "17px",
                  marginTop: "auto",
                  marginBottom: "auto",
                  marginRight: "10px",
                }}
              />
            ) : (
              <></>
            )}
            {submitting ? "Guardando..." : "Guardar"}
          </Button>
        </Col>
      </Row>
    </fieldset>
  );
};

export default ConfigForm;
